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
			total := 0
			for _, n := range rep.Orphans {
				total += n
			}
			fmt.Fprintf(cmd.OutOrStdout(), "imported: %d tables; discogs coverage: %v; orphan total: %d\n", len(rep.Counts), rep.DiscogsCoverage, total)
			writeReportDetail(cmd.OutOrStdout(), rep)
			return nil
		},
	}
	c.Flags().StringVar(&dumpDir, "dump-dir", "", "directory containing mbdump*.tar.bz2 and canonical_*.csv")
	c.MarkFlagRequired("dump-dir")
	return c
}
