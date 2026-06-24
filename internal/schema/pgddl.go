// internal/schema/pgddl.go
//
// A small, regex-free PostgreSQL-DDL scanner. It is paren-depth aware and
// single-quote aware so it survives the real MusicBrainz CreateTables.sql:
// comments containing parens, numeric(p,s), nested inline CHECK constraints,
// and string literals that contain '--', ';' or ')'.
package schema

import "strings"

type rawCol struct {
	Name   string
	PGType string
}

var constraintKW = map[string]bool{
	"CONSTRAINT": true, "CHECK": true, "PRIMARY": true,
	"FOREIGN": true, "UNIQUE": true, "EXCLUDE": true,
}

// stripComments removes `-- line` and `/* block */` comments while preserving
// the contents of single-quoted string literals (a `--` or `/*` inside a quote
// is data, not a comment). PostgreSQL escapes a quote inside a literal by
// doubling it (”) — that is handled naturally here since the doubled quote
// just toggles in-quote off then on again.
func stripComments(sql string) string {
	var b strings.Builder
	b.Grow(len(sql))
	inQuote := false
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inQuote {
			b.WriteByte(c)
			if c == '\'' {
				inQuote = false
			}
			continue
		}
		switch {
		case c == '\'':
			inQuote = true
			b.WriteByte(c)
		case c == '-' && i+1 < len(sql) && sql[i+1] == '-':
			// line comment: skip to end of line (keep the newline)
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			if i < len(sql) {
				b.WriteByte('\n')
			}
		case c == '/' && i+1 < len(sql) && sql[i+1] == '*':
			// block comment: skip to closing */
			i += 2
			for i+1 < len(sql) && !(sql[i] == '*' && sql[i+1] == '/') {
				i++
			}
			i++ // land on '/', loop's i++ advances past it
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// splitStatements splits SQL into statements on `;` at paren-depth 0, honoring
// single-quoted literals. Empty (whitespace-only) statements are dropped.
func splitStatements(sql string) []string {
	var stmts []string
	depth := 0
	inQuote := false
	start := 0
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inQuote {
			if c == '\'' {
				inQuote = false
			}
			continue
		}
		switch c {
		case '\'':
			inQuote = true
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				if s := strings.TrimSpace(sql[start:i]); s != "" {
					stmts = append(stmts, s)
				}
				start = i + 1
			}
		}
	}
	if s := strings.TrimSpace(sql[start:]); s != "" {
		stmts = append(stmts, s)
	}
	return stmts
}

// splitTopCommas splits a parenthesized table body on commas at paren-depth 0,
// honoring single-quoted literals, so numeric(p,s) and CHECK (... IN (a,b))
// stay within a single item.
func splitTopCommas(body string) []string {
	var items []string
	depth := 0
	inQuote := false
	start := 0
	for i := 0; i < len(body); i++ {
		c := body[i]
		if inQuote {
			if c == '\'' {
				inQuote = false
			}
			continue
		}
		switch c {
		case '\'':
			inQuote = true
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				items = append(items, body[start:i])
				start = i + 1
			}
		}
	}
	items = append(items, body[start:])
	return items
}

// parenBody returns the substring between the first top-level '(' and its
// matching ')' (quote aware). It returns the body, the index just past the
// closing ')', and ok=false if no balanced parens are found.
func parenBody(s string) (body string, rest string, ok bool) {
	depth := 0
	inQuote := false
	open := -1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			if c == '\'' {
				inQuote = false
			}
			continue
		}
		switch c {
		case '\'':
			inQuote = true
		case '(':
			if depth == 0 {
				open = i
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				return s[open+1 : i], s[i+1:], true
			}
		}
	}
	return "", "", false
}

// ParseTables extracts column name/type per table, in declaration order.
func ParseTables(sql string) (map[string][]rawCol, []string, error) {
	tables := map[string][]rawCol{}
	var order []string
	clean := stripComments(sql)
	for _, stmt := range splitStatements(clean) {
		toks := strings.Fields(stmt)
		if len(toks) < 3 || !strings.EqualFold(toks[0], "CREATE") || !strings.EqualFold(toks[1], "TABLE") {
			continue
		}
		// Skip an optional "IF NOT EXISTS" before the table name.
		idx := 2
		if idx+2 < len(toks) &&
			strings.EqualFold(toks[idx], "IF") &&
			strings.EqualFold(toks[idx+1], "NOT") &&
			strings.EqualFold(toks[idx+2], "EXISTS") {
			idx += 3
		}
		if idx >= len(toks) {
			continue
		}
		name := strings.TrimSuffix(toks[idx], "(")

		// Skip PostgreSQL partition sub-tables (PARTITION OF <parent> ...).
		// They inherit their schema from the parent and have no column body
		// SQLite can use; the parent table already covers them.
		if strings.Contains(strings.ToUpper(stmt), "PARTITION OF") {
			continue
		}

		body, _, ok := parenBody(stmt)
		if !ok {
			continue
		}
		var cols []rawCol
		for _, item := range splitTopCommas(body) {
			fields := strings.Fields(strings.TrimSpace(item))
			if len(fields) < 2 {
				continue // blank or a single-token fragment — no column here
			}
			if constraintKW[strings.ToUpper(fields[0])] {
				continue // table-level constraint
			}
			cols = append(cols, rawCol{Name: fields[0], PGType: typeExpr(fields[1:])})
		}
		tables[name] = cols
		order = append(order, name)
	}
	return tables, order, nil
}

// typeExpr greedily joins multi-word types (e.g. TIMESTAMP WITH TIME ZONE, CHARACTER VARYING)
// and stops at the first modifier keyword.
func typeExpr(toks []string) string {
	var out []string
	for _, tk := range toks {
		u := strings.ToUpper(strings.TrimSuffix(tk, ","))
		switch u {
		case "NOT", "NULL", "DEFAULT", "CONSTRAINT", "CHECK", "REFERENCES", "PRIMARY", "UNIQUE":
			return strings.Join(out, " ")
		case "WITH", "TIME", "ZONE", "VARYING", "PRECISION", "DOUBLE":
			out = append(out, u) // part of a multi-word type
		default:
			if len(out) > 0 {
				// already captured the base type; a bare extra token is a modifier → stop
				return strings.Join(out, " ")
			}
			out = append(out, strings.TrimSuffix(tk, ","))
		}
	}
	return strings.Join(out, " ")
}

// alterTableTarget returns the table name for an `ALTER TABLE <name> ...`
// statement (already comment-stripped), or "" if the statement is not one.
func alterTableTarget(toks []string) string {
	if len(toks) < 3 || !strings.EqualFold(toks[0], "ALTER") || !strings.EqualFold(toks[1], "TABLE") {
		return ""
	}
	i := 2
	if strings.EqualFold(toks[i], "ONLY") { // ALTER TABLE ONLY <name>
		i++
	}
	if i >= len(toks) {
		return ""
	}
	return toks[i]
}

// keyColumns finds `<kw1> <kw2> ( col, col, ... )` in the statement
// (e.g. PRIMARY KEY (...), FOREIGN KEY (...)) and returns the column list.
// It walks the statement token-by-token (tracking each token's byte offset)
// and, on finding the adjacent keyword pair, extracts the following
// parenthesized list with parenBody so nested parens are handled correctly.
func keyColumns(stmt, kw1, kw2 string) ([]string, bool) {
	type tok struct {
		text string
		end  int // byte offset just past this token in stmt
	}
	var toks []tok
	pos := 0
	for _, f := range strings.Fields(stmt) {
		i := strings.Index(stmt[pos:], f)
		pos += i + len(f)
		toks = append(toks, tok{text: f, end: pos})
	}
	for i := 0; i+1 < len(toks); i++ {
		if !strings.EqualFold(toks[i].text, kw1) || !strings.EqualFold(toks[i+1].text, kw2) {
			continue
		}
		body, _, ok := parenBody(stmt[toks[i+1].end:])
		if !ok {
			return nil, false
		}
		var cols []string
		for _, c := range splitTopCommas(body) {
			if t := strings.TrimSpace(c); t != "" {
				cols = append(cols, t)
			}
		}
		return cols, len(cols) > 0
	}
	return nil, false
}

func ParsePrimaryKeys(sql string) map[string][]string {
	out := map[string][]string{}
	clean := stripComments(sql)
	for _, stmt := range splitStatements(clean) {
		toks := strings.Fields(stmt)
		name := alterTableTarget(toks)
		if name == "" {
			continue
		}
		if cols, ok := keyColumns(stmt, "PRIMARY", "KEY"); ok {
			out[name] = cols
		}
	}
	return out
}

func ParseForeignKeys(sql string) [][2]string {
	var out [][2]string
	clean := stripComments(sql)
	for _, stmt := range splitStatements(clean) {
		toks := strings.Fields(stmt)
		name := alterTableTarget(toks)
		if name == "" {
			continue
		}
		if cols, ok := keyColumns(stmt, "FOREIGN", "KEY"); ok {
			out = append(out, [2]string{name, cols[0]})
		}
	}
	return out
}
