// internal/schema/types.go
package schema

import "strings"

type Column struct {
	Name string `json:"name"`
	Bool bool   `json:"bool,omitempty"`
}

type Manifest struct {
	SchemaSequence int                 `json:"schema_sequence"`
	Tables         map[string][]Column `json:"tables"`
}

// MapType maps a PostgreSQL type expression to a SQLite type and whether it is boolean.
func MapType(pg string) (string, bool) {
	u := strings.ToUpper(strings.TrimSpace(pg))
	if strings.HasSuffix(u, "[]") {
		return "TEXT", false // arrays stored as the {..} literal
	}
	if i := strings.IndexByte(u, '('); i >= 0 {
		u = strings.TrimSpace(u[:i]) // drop length: VARCHAR(255) -> VARCHAR
	}
	switch u {
	case "SERIAL", "BIGSERIAL", "SMALLSERIAL", "INTEGER", "INT", "INT4", "SMALLINT", "INT2", "BIGINT", "INT8":
		return "INTEGER", false
	case "BOOLEAN", "BOOL":
		return "INTEGER", true
	case "NUMERIC", "DECIMAL", "REAL", "DOUBLE PRECISION":
		return "REAL", false
	case "UUID", "TEXT", "VARCHAR", "CHARACTER VARYING", "CHAR", "CHARACTER",
		"TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMPTZ", "DATE", "TIME",
		"POINT", "CUBE", "BYTEA", "JSON", "JSONB", "INTERVAL":
		return "TEXT", false
	default:
		return "TEXT", false // custom ENUM/domain types → TEXT
	}
}
