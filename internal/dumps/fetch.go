package dumps

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// canonicalDirRe matches directory entries in the autoindex. The trailing
// slash anchors the match to directories (never a bare filename like the
// .tar.zst archive); strip it before comparing/returning.
var canonicalDirRe = regexp.MustCompile(`musicbrainz-canonical-dump-\d{8}-\d{6}/`)

type Client struct {
	HTTP *http.Client
	Base string // e.g. https://data.metabrainz.org/pub/musicbrainz/data
}

func (c *Client) httpc() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func (c *Client) ResolveLatest(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Base+"/fullexport/LATEST", nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpc().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LATEST: status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func (c *Client) Download(ctx context.Context, url, dest string) error {
	var start int64
	if fi, err := os.Stat(dest); err == nil {
		start = fi.Size()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if start > 0 {
		// If dest is already the full size, the server answers 416; we surface that as an
		// error (acceptable for monthly use — remove the partial/complete file to re-fetch).
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", start))
	}
	resp, err := c.httpc().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	flags := os.O_CREATE | os.O_WRONLY
	if resp.StatusCode == http.StatusPartialContent {
		flags |= os.O_APPEND
	} else if resp.StatusCode == http.StatusOK {
		flags |= os.O_TRUNC
	} else {
		return fmt.Errorf("download %s: status %d", url, resp.StatusCode)
	}
	f, err := os.OpenFile(dest, flags, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func VerifySHA256(dir, sumsFile string, files []string) error {
	want := map[string]string{}
	f, err := os.Open(filepath.Join(dir, sumsFile))
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 2 {
			want[fields[1]] = fields[0]
		}
	}
	for _, name := range files {
		exp, ok := want[name]
		if !ok {
			return fmt.Errorf("%s: no checksum in %s", name, sumsFile)
		}
		got, err := sha256File(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if !strings.EqualFold(got, exp) {
			return fmt.Errorf("%s: sha256 mismatch (got %s want %s)", name, got, exp)
		}
	}
	return nil
}

func (c *Client) canonicalRoot() string {
	return strings.TrimSuffix(c.Base, "/data") + "/canonical_data"
}

// ResolveLatestCanonical fetches the canonical-data directory listing and
// returns the name of the newest dump directory (lexical maximum of all
// musicbrainz-canonical-dump-YYYYMMDD-HHMMSS entries).
func (c *Client) ResolveLatestCanonical(ctx context.Context) (string, error) {
	url := c.canonicalRoot() + "/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpc().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("canonical index: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	matches := canonicalDirRe.FindAllString(string(body), -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("canonical index: no dump directories found")
	}
	latest := strings.TrimSuffix(matches[0], "/")
	for _, m := range matches[1:] {
		if d := strings.TrimSuffix(m, "/"); d > latest {
			latest = d
		}
	}
	return latest, nil
}

// DownloadCanonical downloads the tar.zst archive for the given dump directory
// into dest, verifies its SHA256 checksum, and returns the local path to the
// archive.
func (c *Client) DownloadCanonical(ctx context.Context, dir, dest string) (string, error) {
	root := c.canonicalRoot()
	archiveName := dir + ".tar.zst"
	sumName := archiveName + ".sha256"

	archivePath := filepath.Join(dest, archiveName)
	sumPath := filepath.Join(dest, sumName)

	if err := c.Download(ctx, root+"/"+dir+"/"+archiveName, archivePath); err != nil {
		return "", fmt.Errorf("download archive: %w", err)
	}
	if err := c.Download(ctx, root+"/"+dir+"/"+sumName, sumPath); err != nil {
		return "", fmt.Errorf("download checksum: %w", err)
	}
	if err := VerifySHA256(dest, sumName, []string{archiveName}); err != nil {
		return "", fmt.Errorf("verify: %w", err)
	}
	return archivePath, nil
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
