package encodingex

import (
	"io"
)

type ignoreReader struct {
	r io.Reader
	byteMap map[byte]bool
}

func NewIgnoreReader(r io.Reader, ignoreBytes []byte) io.Reader {
	byteMap := make(map[byte]bool)
	for _, b := range ignoreBytes {
		byteMap[b] = true
	}
	return &ignoreReader{r, byteMap}
}

func (r *ignoreReader) Read(p []byte) (n int, err error) {
	buf := make([]byte, len(p), cap(p))
	var bufn int
	bufn, err = r.r.Read(buf)
	if err != nil {
		return
	}
	n = 0
	for _, b := range buf[:bufn] {
		if _, ok := r.byteMap[b]; !ok {
			p[n] = b
			n++
		}
	}
	return
}

