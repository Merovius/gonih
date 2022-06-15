// Package lines wraps a Reader to map file offsets to line information.
//
// The canonical use case are parsers like encoding/xml and encoding/json. They
// report errors and token positions as byte-offsets in the input. When
// reporting errors to the user, these offsets are mostly useless. This package
// allows translating them into line/colum numbers. See the examples for how to
// do that.
package lines

import (
	"bytes"
	"io"
	"sort"
	"sync"
)

// A Reader wraps an io.Reader and keeps track of line information read through
// it. It is safe for concurrent use.
type Reader struct {
	mu    sync.RWMutex
	r     io.Reader
	offs  int64
	lines []int64
}

// NewReader wrap r to keep track of line information.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// Read passes all calls through to the underlying io.Reader, recording line
// endings encountered in the streamed data.
func (r *Reader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, err = r.r.Read(p)
	for i := 0; i < n; {
		j := bytes.IndexByte(p[i:n], '\n')
		if j < 0 {
			break
		}
		i += j + 1
		r.lines = append(r.lines, r.offs+int64(i))
	}
	r.offs += int64(n)
	return n, err
}

// Size is the number of bytes read so far. Position information is only
// accurate for offsets less than Size.
func (r *Reader) Size() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.offs
}

// Position returns the line and column of given offset (in bytes). Lines and
// columns are numbered starting with 1. End-of-line markers are counted as
// part of the line preceding them. "\n" is used as an end-of-line marker,
// which also covers systems where "\r\n" is canonically used.
//
// The returned information is only accurate if offset is less than Size.
func (r *Reader) Position(offs int64) (line, column int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if offs < 0 {
		panic("Position called with negative offset")
	}
	i := sort.Search(len(r.lines), func(i int) bool {
		return offs < r.lines[i]
	})
	if i == 0 {
		return 1, int(offs + 1)
	}
	return i + 1, int(offs - r.lines[i-1] + 1)
}

// Line is like Position, but only returns the line.
func (r *Reader) Line(offs int64) int {
	l, _ := r.Position(offs)
	return l
}

// Column is like Position, but only returns the column.
func (r *Reader) Column(offs int64) int {
	_, c := r.Position(offs)
	return c
}

// Reset the recorded position information and continue reading from nr.
func (r *Reader) Reset(nr io.Reader) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.r = nr
	r.offs = 0
	r.lines = r.lines[:0]
}
