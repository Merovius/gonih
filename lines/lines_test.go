package lines

import (
	"io"
	"log"
	"strings"
	"testing"
)

func TestLines(t *testing.T) {
	tcs := []struct {
		in       string
		offs     int64
		wantLine int
		wantCol  int
	}{
		{"foo\nbar\nbaz", 0, 1, 1},
		{"foo\nbar\nbaz", 3, 1, 4},
		{"foo\nbar\nbaz", 4, 2, 1},
		{"foo\nbar\nbaz", 5, 2, 2},
		{"foo\nbar\nbaz", 7, 2, 4},
		{"foo\nbar\nbaz", 8, 3, 1},
		{"foo\nbar\nbaz", 20, 3, 13},
		{"foo\nbar\nbaz\n", 20, 4, 9},
	}
	for _, tc := range tcs {
		r := NewReader(strings.NewReader(tc.in))
		io.Copy(io.Discard, r)
		l, c := r.Position(tc.offs)
		if l != tc.wantLine || c != tc.wantCol {
			log.Printf("Position(%d, %q) = (%d, %d), want (%d, %d)", tc.offs, tc.in, l, c, tc.wantLine, tc.wantCol)
		}
	}
}
