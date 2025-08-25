package str

import (
	"bufio"
	"io"
	"iter"
)

func IterateLines(reader io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if !yield(line) {
				return
			}
		}
	}
}
