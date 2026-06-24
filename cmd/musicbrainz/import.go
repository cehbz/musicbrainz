package main

import (
	"fmt"

	"github.com/cehbz/musicbrainz/internal/store"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	var dumpDir string
	c := &cobra.Command{
		Use:   "import",
		Short: "build a fresh SQLite DB from already-fetched dumps",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dataRoot, _ := cmd.Flags().GetString("data-root")
			rep, err := store.RunImport(dataRoot, dumpDir)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "imported: %d tables; discogs coverage: %v\n", len(rep.Counts), rep.DiscogsCoverage)
			if len(rep.Canonical) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "WARNING: canonical tables are empty\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "canonical: %v\n", rep.Canonical)
			}
			if len(rep.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "skipped %d entries with no matching table:\n", len(rep.Skipped))
				for _, s := range rep.Skipped {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", s)
				}
			}
			if len(rep.Malformed) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "malformed rows (skipped): %v\n", rep.Malformed)
			}
			return nil
		},
	}
	c.Flags().StringVar(&dumpDir, "dump-dir", "", "directory containing mbdump*.tar.bz2 and canonical_*.csv")
	c.MarkFlagRequired("dump-dir")
	return c
}
