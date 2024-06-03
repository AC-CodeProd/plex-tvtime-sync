package share

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	noAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	multipleSpacesRegex = regexp.MustCompile(`\s+`)
)

func CloseFile(f io.ReadCloser) {
	if f != nil {
		f.Close()
	}
}

func ClearString(s string) string {
	cleaned := noAlphanumericRegex.ReplaceAllString(s, "")

	cleaned = strings.TrimSpace(cleaned)
	cleaned = multipleSpacesRegex.ReplaceAllString(cleaned, " ")

	return cleaned
}

func EnsureDir(path string) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
