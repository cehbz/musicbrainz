// internal/schema/translate_test.go
package schema

import "testing"

func TestMapType(t *testing.T) {
	cases := []struct {
		in   string
		want string
		bool bool
	}{
		{"SERIAL", "INTEGER", false},
		{"INTEGER", "INTEGER", false},
		{"SMALLINT", "INTEGER", false},
		{"BIGINT", "INTEGER", false},
		{"BOOLEAN", "INTEGER", true},
		{"UUID", "TEXT", false},
		{"TEXT", "TEXT", false},
		{"VARCHAR", "TEXT", false},
		{"VARCHAR(255)", "TEXT", false},
		{"CHARACTER VARYING", "TEXT", false},
		{"TIMESTAMP WITH TIME ZONE", "TEXT", false},
		{"DATE", "TEXT", false},
		{"INTEGER[]", "TEXT", false},
		{"POINT", "TEXT", false},
		{"CUBE", "TEXT", false},
		{"SOME_ENUM_TYPE", "TEXT", false}, // unknown custom type → TEXT
	}
	for _, c := range cases {
		got, b := MapType(c.in)
		if got != c.want || b != c.bool {
			t.Errorf("MapType(%q) = (%q,%v), want (%q,%v)", c.in, got, b, c.want, c.bool)
		}
	}
}

func TestTranslate(t *testing.T) {
	createTables := `
CREATE TABLE artist ( -- replicate (verbose)
    id                  SERIAL,
    gid                 UUID NOT NULL,
    name                VARCHAR NOT NULL,
    ended               BOOLEAN NOT NULL DEFAULT FALSE
        CONSTRAINT artist_ended_check CHECK (ended IS NOT NULL),
    last_updated        TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE artist_credit_name (
    artist_credit       INTEGER NOT NULL,
    position            SMALLINT NOT NULL,
    artist              INTEGER NOT NULL,
    name                VARCHAR NOT NULL
);
`
	primaryKeys := `
ALTER TABLE artist ADD CONSTRAINT artist_pkey PRIMARY KEY (id);
ALTER TABLE artist_credit_name ADD CONSTRAINT artist_credit_name_pkey PRIMARY KEY (artist_credit, position);
`
	fks := `ALTER TABLE artist_credit_name ADD CONSTRAINT x FOREIGN KEY (artist) REFERENCES artist(id);`

	schemaSQL, indexesSQL, m, err := Translate(createTables, "", primaryKeys, fks, 31)
	if err != nil {
		t.Fatal(err)
	}
	// single-int PK id becomes rowid alias
	if !contains(schemaSQL, "id INTEGER PRIMARY KEY") {
		t.Errorf("artist.id not a rowid PK:\n%s", schemaSQL)
	}
	// boolean mapped to INTEGER; constraint line dropped
	if !contains(schemaSQL, "ended INTEGER") || contains(schemaSQL, "CHECK") {
		t.Errorf("boolean/constraint handling wrong:\n%s", schemaSQL)
	}
	// enriched discogs column appended to artist
	if !contains(schemaSQL, "discogs_artist_id INTEGER") {
		t.Errorf("missing enriched discogs column:\n%s", schemaSQL)
	}
	// composite PK -> post-load unique index, gid -> unique index, FK -> index
	if !contains(indexesSQL, "ON artist_credit_name(artist_credit, position)") {
		t.Errorf("composite PK index missing:\n%s", indexesSQL)
	}
	if !contains(indexesSQL, "ON artist(gid)") {
		t.Errorf("gid unique index missing:\n%s", indexesSQL)
	}
	if !contains(indexesSQL, "ON artist_credit_name(artist)") {
		t.Errorf("FK index missing:\n%s", indexesSQL)
	}
	// manifest preserves column order and bool flag
	cols := m.Tables["artist"]
	if len(cols) != 5 || cols[0].Name != "id" || cols[3].Name != "ended" || !cols[3].Bool {
		t.Errorf("artist manifest wrong: %+v", cols)
	}
	if m.SchemaSequence != 31 {
		t.Errorf("sequence = %d, want 31", m.SchemaSequence)
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && indexOf(s, sub) >= 0 }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
