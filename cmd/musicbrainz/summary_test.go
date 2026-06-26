package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cehbz/musicbrainz/internal/store"
)

func TestWriteReportDetailSurfacesSkippedAndMalformed(t *testing.T) {
	rep := store.Report{
		Canonical: map[string]int{"canonical_musicbrainz_data": 5},
		Skipped:   []string{"mbdump-edit", "mbdump-stats"},
		Malformed: map[string]int{"artist": 3},
	}
	var out bytes.Buffer
	writeReportDetail(&out, rep)
	s := out.String()
	for _, want := range []string{"skipped 2 entries", "mbdump-edit", "mbdump-stats", "malformed", "artist", "3"} {
		if !strings.Contains(s, want) {
			t.Errorf("output missing %q; got:\n%s", want, s)
		}
	}
}
