// internal/dumps/fetch_test.go
package dumps

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveLatestCanonical(t *testing.T) {
	const indexHTML = `<!DOCTYPE HTML>
<html><head><title>Index of /pub/musicbrainz/canonical_data/</title></head>
<body>
<a href="../">../</a>
<a href="musicbrainz-canonical-dump-20260603-080003/">musicbrainz-canonical-dump-20260603-080003/</a>
<a href="musicbrainz-canonical-dump-20260617-080003/">musicbrainz-canonical-dump-20260617-080003/</a>
</body></html>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/canonical_data/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(indexHTML))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c := &Client{HTTP: srv.Client(), Base: srv.URL + "/data"}
	got, err := c.ResolveLatestCanonical(context.Background())
	if err != nil {
		t.Fatalf("ResolveLatestCanonical error: %v", err)
	}
	const want = "musicbrainz-canonical-dump-20260617-080003"
	if got != want {
		t.Fatalf("ResolveLatestCanonical = %q, want %q", got, want)
	}
}

func TestDownloadCanonical(t *testing.T) {
	const dir = "musicbrainz-canonical-dump-20260617-080003"
	const archiveBody = "fake-tar-zst-data"
	sum := sha256.Sum256([]byte(archiveBody))
	sumLine := hex.EncodeToString(sum[:]) + "  " + dir + ".tar.zst\n"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/canonical_data/" + dir + "/" + dir + ".tar.zst":
			w.Write([]byte(archiveBody))
		case "/canonical_data/" + dir + "/" + dir + ".tar.zst.sha256":
			w.Write([]byte(sumLine))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := &Client{HTTP: srv.Client(), Base: srv.URL + "/data"}
	tmp := t.TempDir()
	tarPath, err := c.DownloadCanonical(context.Background(), dir, tmp)
	if err != nil {
		t.Fatalf("DownloadCanonical error: %v", err)
	}
	got, err := os.ReadFile(tarPath)
	if err != nil {
		t.Fatalf("read archive: %v", err)
	}
	if string(got) != archiveBody {
		t.Fatalf("archive body = %q, want %q", got, archiveBody)
	}
}

func TestDownloadCanonicalChecksumMismatch(t *testing.T) {
	const dir = "musicbrainz-canonical-dump-20260617-080003"
	const archiveBody = "fake-tar-zst-data"
	// Deliberately wrong hash: sum of different bytes.
	badSum := sha256.Sum256([]byte("not-the-archive"))
	sumLine := hex.EncodeToString(badSum[:]) + "  " + dir + ".tar.zst\n"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/canonical_data/" + dir + "/" + dir + ".tar.zst":
			w.Write([]byte(archiveBody))
		case "/canonical_data/" + dir + "/" + dir + ".tar.zst.sha256":
			w.Write([]byte(sumLine))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := &Client{HTTP: srv.Client(), Base: srv.URL + "/data"}
	if _, err := c.DownloadCanonical(context.Background(), dir, t.TempDir()); err == nil {
		t.Fatal("DownloadCanonical: expected error on checksum mismatch, got nil")
	}
}

func TestResolveLatestAndVerify(t *testing.T) {
	body := "1\tfoo\n"
	sum := sha256.Sum256([]byte(body))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fullexport/LATEST":
			w.Write([]byte("20260620-002052\n"))
		case "/fullexport/20260620-002052/mbdump.tar.bz2":
			w.Write([]byte(body))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := &Client{HTTP: srv.Client(), Base: srv.URL}
	dir, err := c.ResolveLatest(context.Background())
	if err != nil || dir != "20260620-002052" {
		t.Fatalf("ResolveLatest = %q, %v", dir, err)
	}
	tmp := t.TempDir()
	dest := filepath.Join(tmp, "mbdump.tar.bz2")
	if err := c.Download(context.Background(), srv.URL+"/fullexport/20260620-002052/mbdump.tar.bz2", dest); err != nil {
		t.Fatal(err)
	}
	// write a matching SHA256SUMS and verify
	if err := os.WriteFile(filepath.Join(tmp, "SHA256SUMS"), []byte(hex.EncodeToString(sum[:])+"  mbdump.tar.bz2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := VerifySHA256(tmp, "SHA256SUMS", []string{"mbdump.tar.bz2"}); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestDownloadResume(t *testing.T) {
	const full = "0123456789ABCDEFGHIJ" // 20 bytes
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rng := r.Header.Get("Range")
		if rng == "" {
			w.Write([]byte(full))
			return
		}
		var start int
		fmt.Sscanf(rng, "bytes=%d-", &start)
		if start >= len(full) {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(full)-1, len(full)))
		w.WriteHeader(http.StatusPartialContent)
		w.Write([]byte(full[start:]))
	}))
	defer srv.Close()

	c := &Client{HTTP: srv.Client(), Base: srv.URL}
	dest := filepath.Join(t.TempDir(), "f.bin")
	if err := os.WriteFile(dest, []byte(full[:8]), 0o644); err != nil { // pre-seed a partial download
		t.Fatal(err)
	}
	if err := c.Download(context.Background(), srv.URL+"/f.bin", dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != full {
		t.Fatalf("resumed file = %q, want %q", got, full)
	}
}
