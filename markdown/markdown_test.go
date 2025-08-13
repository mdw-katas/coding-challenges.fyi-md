package markdown

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mdw-go/set/v2/set"
	"github.com/mdw-go/testing/v2/assert"
	"github.com/mdw-go/testing/v2/better"
	"github.com/mdw-go/testing/v2/should"
	"github.com/yuin/goldmark"
)

//go:embed testdata
var testFiles embed.FS

func iterateTestsCases(t *testing.T) iter.Seq[string] {
	tests := set.Of[string]()
	err := fs.WalkDir(testFiles, "testdata", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			tests.Add(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		}
		return nil
	})
	assert.So(t, err, better.BeNil)
	return tests.All()
}

func Test(t *testing.T) {
	standard := goldmark.New() // TODO: options and extensions
	for _, testID := range slices.Sorted(iterateTestsCases(t)) {
		t.Run(testID, func(t *testing.T) {
			filename := fmt.Sprintf("testdata/%s", testID)
			input, err := testFiles.ReadFile(filename + ".md")
			assert.So(t, err, better.BeNil)
			expected, err := testFiles.ReadFile(filename + ".html")
			if errors.Is(err, os.ErrNotExist) {
				var buffer bytes.Buffer
				_ = standard.Convert(input, &buffer)
				_ = os.WriteFile(filename+".html", buffer.Bytes(), 0644)
				expected = buffer.Bytes()
			}
			assert.So(t, ConvertToHTML(string(input)), should.Equal, string(expected))
		})
	}
}
