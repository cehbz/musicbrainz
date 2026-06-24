// internal/schema/translate.go
package schema

import (
	"fmt"
	"sort"
	"strings"
)

// discogsCols are advisory cross-DB FK columns appended to the served schema.
var discogsCols = map[string]string{
	"artist":        "discogs_artist_id",
	"label":         "discogs_label_id",
	"release_group": "discogs_master_id",
	"release":       "discogs_release_id",
}

func Translate(createTables, createTypes, primaryKeys, fkConstraints string, seq int) (string, string, Manifest, error) {
	tables, order, err := ParseTables(createTables)
	if err != nil {
		return "", "", Manifest{}, err
	}
	pks := ParsePrimaryKeys(primaryKeys)
	fks := ParseForeignKeys(fkConstraints)

	m := Manifest{SchemaSequence: seq, Tables: map[string][]Column{}}
	var schemaB, idxB strings.Builder

	for _, tname := range order {
		cols := tables[tname]
		pk := pks[tname]
		singleIntPK := len(pk) == 1 // promote to rowid alias if the column is integer

		var defs []string
		var mcols []Column
		for _, c := range cols {
			sqlType, isBool := MapType(c.PGType)
			mcols = append(mcols, Column{Name: c.Name, Bool: isBool})
			def := fmt.Sprintf("  %s %s", c.Name, sqlType)
			if singleIntPK && c.Name == pk[0] && sqlType == "INTEGER" {
				def = fmt.Sprintf("  %s INTEGER PRIMARY KEY", c.Name)
				singleIntPK = false // consumed
			}
			defs = append(defs, def)
		}
		if extra, ok := discogsCols[tname]; ok {
			defs = append(defs, fmt.Sprintf("  %s INTEGER", extra))
		}
		m.Tables[tname] = mcols[:len(cols)] // manifest covers ONLY dump columns (order/bool)
		fmt.Fprintf(&schemaB, "CREATE TABLE %s (\n%s\n);\n", tname, strings.Join(defs, ",\n"))

		// post-load indexes
		if len(pk) > 1 {
			fmt.Fprintf(&idxB, "CREATE UNIQUE INDEX %s_pkey ON %s(%s);\n", tname, tname, strings.Join(pk, ", "))
		} else if len(pk) == 1 {
			// pk became rowid alias only if integer; if not promoted, index it
			if cType, _ := MapType(colType(cols, pk[0])); cType != "INTEGER" {
				fmt.Fprintf(&idxB, "CREATE UNIQUE INDEX %s_pkey ON %s(%s);\n", tname, tname, pk[0])
			}
		}
		if hasCol(cols, "gid") {
			fmt.Fprintf(&idxB, "CREATE UNIQUE INDEX %s_gid ON %s(gid);\n", tname, tname)
		}
		if extra, ok := discogsCols[tname]; ok {
			fmt.Fprintf(&idxB, "CREATE INDEX %s_%s ON %s(%s);\n", tname, extra, tname, extra)
		}
	}

	// FK-join indexes (deterministic order)
	sort.Slice(fks, func(i, j int) bool {
		if fks[i][0] != fks[j][0] {
			return fks[i][0] < fks[j][0]
		}
		return fks[i][1] < fks[j][1]
	})
	for _, fk := range fks {
		fmt.Fprintf(&idxB, "CREATE INDEX %s_%s_idx ON %s(%s);\n", fk[0], fk[1], fk[0], fk[1])
	}
	return schemaB.String(), idxB.String(), m, nil
}

func colType(cols []rawCol, name string) string {
	for _, c := range cols {
		if c.Name == name {
			return c.PGType
		}
	}
	return ""
}
func hasCol(cols []rawCol, name string) bool { return colType(cols, name) != "" }
