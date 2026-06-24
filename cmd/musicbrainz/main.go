package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

func newRootCmd() *cobra.Command {
	root := &cobra.Command{Use: "musicbrainz", SilenceUsage: true}
	root.PersistentFlags().String("data-root", "./musicbrainz-data", "directory holding dumps and the DB")
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "musicbrainz %s\n", Version)
			return nil
		},
	})
	return root
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
