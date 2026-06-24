package main

import (
	"fmt"

	"github.com/cehbz/musicbrainz/internal/store"
	"github.com/spf13/cobra"
)

func newVerifyCmd() *cobra.Command {
	var dbPath string
	c := &cobra.Command{
		Use:   "verify",
		Short: "print row counts and integrity summary for a built DB",
		RunE: func(cmd *cobra.Command, _ []string) error {
			db, err := store.Open(dbPath, store.ModeServe)
			if err != nil {
				return err
			}
			defer db.Close()
			counts, err := db.TableCounts([]string{"artist", "release_group", "release", "recording", "work", "label", "isrc"})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "counts: %v\n", counts)
			return nil
		},
	}
	c.Flags().StringVar(&dbPath, "db", "", "path to the SQLite DB")
	c.MarkFlagRequired("db")
	return c
}
