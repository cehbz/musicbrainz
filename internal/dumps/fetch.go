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
	"strings"
)

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
