package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cehbz/musicbrainz/internal/dumps"
	"github.com/spf13/cobra"
)

const defaultDumpBase = "https://data.metabrainz.org/pub/musicbrainz/data"

// runFetch resolves the export dir (or uses `date`), downloads the core+derived tarballs
// and SHA256SUMS into dest, verifies them, fetches the canonical dump, and returns the
// resolved fullexport dir name.
func runFetch(ctx context.Context, dest, date string) (string, error) {
	cl := &dumps.Client{Base: defaultDumpBase}
	dir := date
	if dir == "" || dir == "LATEST" {
		d, err := cl.ResolveLatest(ctx)
		if err != nil {
			return "", err
		}
		dir = d
	}
	base := cl.Base + "/fullexport/" + dir
	files := []string{"mbdump.tar.bz2", "mbdump-derived.tar.bz2", "SHA256SUMS"}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return "", err
	}
	for _, f := range files {
		if err := cl.Download(ctx, base+"/"+f, filepath.Join(dest, f)); err != nil {
			return "", err
		}
	}
	if err := dumps.VerifySHA256(dest, "SHA256SUMS", files[:2]); err != nil {
		return "", err
	}

	// fetch canonical dump
	cdir, err := cl.ResolveLatestCanonical(ctx)
	if err != nil {
		return "", fmt.Errorf("resolve canonical: %w", err)
	}
	if _, err := cl.DownloadCanonical(ctx, cdir, dest); err != nil {
		return "", fmt.Errorf("download canonical: %w", err)
	}
	fmt.Printf("fetched canonical %s\n", cdir)

	return dir, nil
}

func newFetchCmd() *cobra.Command {
	var dest, date string
	c := &cobra.Command{
		Use:   "fetch",
		Short: "download and verify the MusicBrainz full-export + canonical dumps",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dir, err := runFetch(cmd.Context(), dest, date)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "fetched fullexport %s into %s\n", dir, dest)
			return nil
		},
	}
	c.Flags().StringVar(&dest, "dest", "", "download directory")
	c.Flags().StringVar(&date, "date", "LATEST", "fullexport dir (YYYYMMDD-HHMMSS) or LATEST")
	c.MarkFlagRequired("dest")
	return c
}
