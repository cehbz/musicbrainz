// internal/schema/pgddl_test.go
package schema

import "testing"

// findCol returns the rawCol with the given name, or nil if absent.
func findCol(cols []rawCol, name string) *rawCol {
	for i := range cols {
		if cols[i].Name == name {
			return &cols[i]
		}
	}
	return nil
}

func colNames(cols []rawCol) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}
	return names
}

// Case 1: a comment-with-parens right after the table head must not corrupt
// paren matching.
func TestParseTables_CommentWithParens(t *testing.T) {
	sql := "CREATE TABLE artist ( -- replicate (verbose)\n" +
		"  id SERIAL,\n" +
		"  gid UUID NOT NULL\n" +
		");"
	tables, order, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 1 || order[0] != "artist" {
		t.Fatalf("table order = %v, want [artist]", order)
	}
	got := colNames(tables["artist"])
	if len(got) != 2 || got[0] != "id" || got[1] != "gid" {
		t.Errorf("artist columns = %v, want [id gid]", got)
	}
}

// Case 2: numeric(p,s) — the comma inside parens must NOT split the column.
func TestParseTables_NumericPrecisionScale(t *testing.T) {
	sql := "CREATE TABLE money_thing (\n" +
		"  id SERIAL,\n" +
		"  amount NUMERIC(5,2) NOT NULL\n" +
		");"
	tables, _, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	cols := tables["money_thing"]
	amount := findCol(cols, "amount")
	if amount == nil {
		t.Fatalf("amount column missing; got %v", colNames(cols))
	}
	if len(cols) != 2 {
		t.Errorf("money_thing has %d cols (%v), want 2", len(cols), colNames(cols))
	}
	// MapType already sends NUMERIC -> REAL.
	if sqlType, _ := MapType(amount.PGType); sqlType != "REAL" {
		t.Errorf("amount PGType=%q maps to %q, want REAL", amount.PGType, sqlType)
	}
}

// Case 3: inline CHECK with nested parens + comma must parse as ONE column.
func TestParseTables_InlineCheckNestedParens(t *testing.T) {
	sql := "CREATE TABLE thing (\n" +
		"  id SERIAL,\n" +
		"  status INTEGER NOT NULL CHECK (status IN (0, 1, 2))\n" +
		");"
	tables, _, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	cols := tables["thing"]
	if len(cols) != 2 {
		t.Errorf("thing has %d cols (%v), want 2", len(cols), colNames(cols))
	}
	status := findCol(cols, "status")
	if status == nil {
		t.Fatalf("status column missing; got %v", colNames(cols))
	}
	if sqlType, _ := MapType(status.PGType); sqlType != "INTEGER" {
		t.Errorf("status PGType=%q maps to %q, want INTEGER", status.PGType, sqlType)
	}
}

// Case 3b: a CHECK whose parens (and commas) span multiple lines, including a
// ')' at the start of a continuation line. This is the case a line-based or
// `\n)\s*;` regex parser gets wrong (it spawns phantom columns).
func TestParseTables_MultiLineCheckNestedParens(t *testing.T) {
	sql := "CREATE TABLE thing (\n" +
		"  id SERIAL,\n" +
		"  status INTEGER NOT NULL\n" +
		"    CHECK (status IN\n" +
		"      (0, 1, 2)\n" +
		"    ),\n" +
		"  name VARCHAR NOT NULL\n" +
		");"
	tables, _, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	got := colNames(tables["thing"])
	if len(got) != 3 || got[0] != "id" || got[1] != "status" || got[2] != "name" {
		t.Errorf("thing columns = %v, want [id status name]", got)
	}
}

// Case 4: table-level constraint lines (CHECK / FOREIGN KEY) produce no column.
func TestParseTables_TableLevelConstraintsSkipped(t *testing.T) {
	sql := "CREATE TABLE thing (\n" +
		"  a INTEGER NOT NULL,\n" +
		"  b INTEGER NOT NULL,\n" +
		"  x INTEGER NOT NULL,\n" +
		"  CONSTRAINT foo_chk CHECK (a > (b + 1)),\n" +
		"  FOREIGN KEY (x) REFERENCES y(id)\n" +
		");"
	tables, _, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	got := colNames(tables["thing"])
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "x" {
		t.Errorf("thing columns = %v, want [a b x]", got)
	}
}

// Case 5: CREATE TABLE IF NOT EXISTS must parse the table name correctly.
func TestParseTables_IfNotExists(t *testing.T) {
	sql := "CREATE TABLE IF NOT EXISTS artist (\n" +
		"  id SERIAL,\n" +
		"  gid UUID NOT NULL\n" +
		");"
	tables, order, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 1 || order[0] != "artist" {
		t.Fatalf("table order = %v, want [artist]", order)
	}
	got := colNames(tables["artist"])
	if len(got) != 2 || got[0] != "id" || got[1] != "gid" {
		t.Errorf("artist columns = %v, want [id gid]", got)
	}
}

// Block comments and single-quoted literals containing '--' or ';' must be
// handled without corrupting statement splitting or comment stripping.
func TestParseTables_BlockCommentsAndQuotedLiterals(t *testing.T) {
	sql := "/* a block comment; with a semicolon and (parens) */\n" +
		"CREATE TABLE thing (\n" +
		"  id SERIAL,\n" +
		"  note VARCHAR NOT NULL DEFAULT 'not -- a comment; not ) a paren'\n" +
		");"
	tables, _, err := ParseTables(sql)
	if err != nil {
		t.Fatal(err)
	}
	got := colNames(tables["thing"])
	if len(got) != 2 || got[0] != "id" || got[1] != "note" {
		t.Errorf("thing columns = %v, want [id note]", got)
	}
}

// ParsePrimaryKeys / ParseForeignKeys must still work without regexp.
func TestParsePrimaryKeys_Scanner(t *testing.T) {
	sql := "ALTER TABLE artist ADD CONSTRAINT artist_pkey PRIMARY KEY (id);\n" +
		"ALTER TABLE acn ADD CONSTRAINT acn_pkey PRIMARY KEY (artist_credit, position);\n"
	pks := ParsePrimaryKeys(sql)
	if got := pks["artist"]; len(got) != 1 || got[0] != "id" {
		t.Errorf("artist PK = %v, want [id]", got)
	}
	if got := pks["acn"]; len(got) != 2 || got[0] != "artist_credit" || got[1] != "position" {
		t.Errorf("acn PK = %v, want [artist_credit position]", got)
	}
}

func TestParseForeignKeys_Scanner(t *testing.T) {
	sql := "ALTER TABLE acn ADD CONSTRAINT x FOREIGN KEY (artist) REFERENCES artist(id);\n"
	fks := ParseForeignKeys(sql)
	if len(fks) != 1 || fks[0][0] != "acn" || fks[0][1] != "artist" {
		t.Errorf("FKs = %v, want [[acn artist]]", fks)
	}
}
