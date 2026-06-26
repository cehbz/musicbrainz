package main

import (
	"fmt"
	"io"

	"github.com/cehbz/musicbrainz/internal/store"
)

// writeReportDetail prints the canonical, skipped, and malformed sections of an
// import Report — the no-silent-drops surfacing shared by `import` and `build`.
func writeReportDetail(w io.Writer, rep store.Report) {
	if len(rep.Canonical) == 0 {
		fmt.Fprintf(w, "WARNING: canonical tables are empty\n")
	} else {
		fmt.Fprintf(w, "canonical: %v\n", rep.Canonical)
	}
	if len(rep.Skipped) > 0 {
		fmt.Fprintf(w, "skipped %d entries with no matching table:\n", len(rep.Skipped))
		for _, s := range rep.Skipped {
			fmt.Fprintf(w, "  %s\n", s)
		}
	}
	if len(rep.Malformed) > 0 {
		fmt.Fprintf(w, "malformed rows (skipped): %v\n", rep.Malformed)
	}
}
