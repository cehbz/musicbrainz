package dumps

import (
	"archive/tar"
	"compress/bzip2"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// streamTar is the tar-walking core (shared by tests and StreamTarBz2).
func streamTar(tr *tar.Reader, fn func(name string, content io.Reader) error) error {
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		if err := fn(h.Name, tr); err != nil {
			return err
		}
	}
}

func StreamTarBz2(r io.Reader, fn func(name string, content io.Reader) error) error {
	return streamTar(tar.NewReader(bzip2.NewReader(r)), fn)
}

// readMetaTar reads SCHEMA_SEQUENCE and TIMESTAMP, stopping as soon as both are seen.
func readMetaTar(tr *tar.Reader) (int, string, error) {
	var seq int
	var ts string
	haveSeq, haveTS := false, false
	for !(haveSeq && haveTS) {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, "", err
		}
		switch h.Name {
		case "SCHEMA_SEQUENCE":
			b, err := io.ReadAll(tr)
			if err != nil {
				return 0, "", err
			}
			seq, err = strconv.Atoi(strings.TrimSpace(string(b)))
			if err != nil {
				return 0, "", err
			}
			haveSeq = true
		case "TIMESTAMP":
			b, err := io.ReadAll(tr)
			if err != nil {
				return 0, "", err
			}
			ts = strings.TrimSpace(string(b))
			haveTS = true
		}
	}
	return seq, ts, nil
}

func ReadMeta(path string) (int, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	return readMetaTar(tar.NewReader(bzip2.NewReader(f)))
}

// ZstdReader returns a streaming decompressor over a zstd stream.
// Close the returned reader to release the decoder's resources.
func ZstdReader(r io.Reader) (io.ReadCloser, error) {
	zr, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	return zr.IOReadCloser(), nil
}
