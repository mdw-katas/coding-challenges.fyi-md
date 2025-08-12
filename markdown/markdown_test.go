package markdown

import (
	"embed"
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mdw-go/set/v2/set"
	"github.com/mdw-go/testing/v2/assert"
	"github.com/mdw-go/testing/v2/better"
	"github.com/mdw-go/testing/v2/should"
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
	for _, testID := range slices.Sorted(iterateTestsCases(t)) {
		t.Run(testID, func(t *testing.T) {
			filename := fmt.Sprintf("testdata/%s", testID)
			expected, err := testFiles.ReadFile(filename + ".html")
			assert.So(t, err, better.BeNil)
			input, err := testFiles.ReadFile(filename + ".md")
			assert.So(t, err, better.BeNil)
			assert.So(t, ConvertToHTML(string(input)), should.Equal, string(expected))
		})
	}
}
