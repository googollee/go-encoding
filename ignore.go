package encodingex

import (
	"io"
)

type ignoreReader struct {
	r       io.Reader
	byteMap map[byte]bool
}

// Create a Reader from r and ignore all ignoreBytes.
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

type ignoreWriter struct {
	w       io.Writer
	byteMap map[byte]bool
}

// Create a Writer from w and ignore all ignoreBytes.
func NewIgnoreWriter(w io.Writer, ignoreBytes []byte) io.Writer {
	byteMap := make(map[byte]bool)
	for _, b := range ignoreBytes {
		byteMap[b] = true
	}
	return &ignoreWriter{w, byteMap}
}

func (w *ignoreWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p), cap(p))
	i := 0
	for _, b := range p {
		if _, ok := w.byteMap[b]; !ok {
			buf[i] = b
			i++
		}
	}
	n, err = w.w.Write(buf[:i])
	if err == nil {
		n = len(p)
	}
	return
}
