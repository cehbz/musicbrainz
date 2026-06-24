package dumps

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
)

// helper: build an in-memory tar.
func writeTar(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, body := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	return buf.Bytes()
}

func TestStreamTarWalksEntries(t *testing.T) {
	tarBytes := writeTar(t, map[string]string{
		"SCHEMA_SEQUENCE": "31\n",
		"mbdump/artist":   "1\tfoo\n2\tbar\n",
	})
	seen := map[string]string{}
	err := streamTar(tar.NewReader(bytes.NewReader(tarBytes)), func(name string, r io.Reader) error {
		b, _ := io.ReadAll(r)
		seen[name] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen["mbdump/artist"] != "1\tfoo\n2\tbar\n" || seen["SCHEMA_SEQUENCE"] != "31\n" {
		t.Fatalf("entries not streamed: %v", seen)
	}
}

func TestZstdReader(t *testing.T) {
	var zb bytes.Buffer
	w, err := zstd.NewWriter(&zb)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("hello,world\n")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	zr, err := ZstdReader(bytes.NewReader(zb.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(zr)
	if string(got) != "hello,world\n" {
		t.Fatalf("zstd roundtrip = %q", got)
	}
}

func TestReadMeta(t *testing.T) {
	tarBytes := writeTar(t, map[string]string{"SCHEMA_SEQUENCE": "31\n", "TIMESTAMP": "2026-06-20 00:20:52\n"})
	seq, ts, err := readMetaTar(tar.NewReader(bytes.NewReader(tarBytes)))
	if err != nil {
		t.Fatal(err)
	}
	if seq != 31 || ts == "" {
		t.Fatalf("ReadMeta = (%d,%q)", seq, ts)
	}
}

func TestStreamTarZst(t *testing.T) {
	files := map[string]string{
		"musicbrainz-canonical-dump-20260617-080003/TIMESTAMP":                                  "20260617-080003\n",
		"musicbrainz-canonical-dump-20260617-080003/canonical/canonical_recording_redirect.csv": "gid,canonical_recording\nfoo,bar\n",
	}
	tarBytes := writeTar(t, files)

	var zb bytes.Buffer
	w, err := zstd.NewWriter(&zb)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(tarBytes); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	seen := map[string]string{}
	err = StreamTarZst(bytes.NewReader(zb.Bytes()), func(name string, r io.Reader) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		seen[name] = string(b)
		return nil
	})
	if err != nil {
		t.Fatalf("StreamTarZst error: %v", err)
	}
	for name, want := range files {
		if got := seen[name]; got != want {
			t.Errorf("entry %q = %q, want %q", name, got, want)
		}
	}
	if len(seen) != len(files) {
		t.Errorf("got %d entries, want %d", len(seen), len(files))
	}
}
