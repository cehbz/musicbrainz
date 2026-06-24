package main

import (
	"os"
	"strings"
	"testing"
)

func TestGeneratedSchemaSane(t *testing.T) {
	b, err := os.ReadFile("../../internal/store/schema_seq31.sql")
	if err != nil {
		t.Fatalf("generated schema missing — run gen-schema: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		"CREATE TABLE artist (",
		"id INTEGER PRIMARY KEY",
		"gid TEXT",
		"discogs_artist_id INTEGER",
		"CREATE TABLE recording (",
		"CREATE TABLE l_recording_work (",
		"CREATE TABLE isrc (",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("generated schema missing %q", want)
		}
	}
}
