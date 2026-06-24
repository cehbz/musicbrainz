package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cehbz/musicbrainz/internal/schema"
	"github.com/spf13/cobra"
)

func newGenSchemaCmd() *cobra.Command {
	var seq int
	var src, out string
	c := &cobra.Command{
		Use:   "gen-schema",
		Short: "translate upstream MusicBrainz DDL into the embedded SQLite schema (dev)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			read := func(name string) (string, error) {
				b, err := os.ReadFile(filepath.Join(src, name))
				if err != nil {
					if os.IsNotExist(err) {
						return "", nil // optional file → empty input is fine
					}
					return "", err
				}
				return string(b), nil
			}
			createTables, err := read("CreateTables.sql")
			if err != nil {
				return err
			}
			createTypes, err := read("CreateTypes.sql")
			if err != nil {
				return err
			}
			primaryKeys, err := read("CreatePrimaryKeys.sql")
			if err != nil {
				return err
			}
			fkConstraints, err := read("CreateFKConstraints.sql")
			if err != nil {
				return err
			}
			schemaSQL, idxSQL, m, err := schema.Translate(
				createTables, createTypes, primaryKeys, fkConstraints, seq)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(out, 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(out, fmt.Sprintf("schema_seq%d.sql", seq)), []byte(schemaSQL), 0o644); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(out, fmt.Sprintf("indexes_seq%d.sql", seq)), []byte(idxSQL), 0o644); err != nil {
				return err
			}
			mj, _ := json.MarshalIndent(m, "", "  ")
			return os.WriteFile(filepath.Join(out, fmt.Sprintf("manifest_seq%d.json", seq)), mj, 0o644)
		},
	}
	c.Flags().IntVar(&seq, "seq", 31, "schema sequence")
	c.Flags().StringVar(&src, "src", "internal/schema/upstream/seq31", "directory of upstream .sql files")
	c.Flags().StringVar(&out, "out", "internal/store", "output directory")
	return c
}
