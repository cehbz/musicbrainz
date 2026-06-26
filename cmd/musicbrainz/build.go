package main

import (
	"fmt"
	"path/filepath"

	"github.com/cehbz/musicbrainz/internal/store"
	"github.com/spf13/cobra"
)

func newBuildCmd() *cobra.Command {
	var date string
	c := &cobra.Command{
		Use:   "build",
		Short: "fetch then import (end-to-end)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dataRoot, _ := cmd.Flags().GetString("data-root")
			dumpDir := filepath.Join(dataRoot, "dumps")
			if _, err := runFetch(cmd.Context(), dumpDir, date); err != nil {
				return err
			}
			rep, err := store.RunImport(dataRoot, dumpDir)
			if err != nil {
				return err
			}
			total := 0
			for _, n := range rep.Orphans {
				total += n
			}
			fmt.Fprintf(cmd.OutOrStdout(), "built: %d tables, discogs coverage %v, orphan total: %d\n", len(rep.Counts), rep.DiscogsCoverage, total)
			writeReportDetail(cmd.OutOrStdout(), rep)
			return nil
		},
	}
	c.Flags().StringVar(&date, "date", "LATEST", "fullexport dir (YYYYMMDD-HHMMSS) or LATEST")
	return c
}
