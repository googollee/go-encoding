package encodingex

import (
	"fmt"
	"io"
)

type IconvReadCloser struct {
	iconv   Iconv
	r       io.Reader
	bufSize int
	inBuf   []byte
	inEnd   int
	outBuf  []byte
	outEnd  int
}

func NewIconvReadCloser(r io.Reader, tocode, fromcode string) (*IconvReadCloser, error) {
	return NewIconvReadCloserBufferSize(r, 1024, tocode, fromcode)
}

func NewIconvReadCloserBufferSize(r io.Reader, bufSize int, tocode, fromcode string) (*IconvReadCloser, error) {
	iconv, err := NewIconv(tocode, fromcode)
	if err != nil {
		return nil, err
	}
	return &IconvReadCloser{
		iconv:   iconv,
		r:       r,
		bufSize: bufSize,
		inBuf:   make([]byte, bufSize, bufSize),
		inEnd:   0,
		outBuf:  make([]byte, bufSize, bufSize),
		outEnd:  0,
	}, nil
}

func (i *IconvReadCloser) Close() error {
	return i.iconv.Close()
}

func (i *IconvReadCloser) Read(p []byte) (int, error) {
	err := i.fillBuffer()
	if i.outEnd == 0 {
		return 0, err
	}
	m := len(p)
	if m > i.outEnd {
		m = i.outEnd
	}
	for j := 0; j < m; j++ {
		p[j] = i.outBuf[j]
	}
	i.moveOutBuffer(m)
	return m, nil
}

func (i *IconvReadCloser) fillBuffer() error {
	if i.outEnd == i.bufSize {
		return nil
	}
	if i.inEnd < i.bufSize {
		n, err := i.r.Read(i.inBuf[i.inEnd:])
		if err != nil && i.inEnd == 0 {
			return err
		}
		i.inEnd += n
	}
	inbuf := i.inBuf[:i.inEnd]
	outbuf := i.outBuf[i.outEnd:]
	inlen, outlen, err := i.iconv.Conv(inbuf, outbuf)
	i.moveOrgBuffer(inlen)
	i.outEnd += outlen
	return err
}

func (i *IconvReadCloser) moveOutBuffer(end int) {
	if end == 0 {
		return
	}
	t := 0
	for f := end; f < i.outEnd; f++ {
		i.outBuf[t] = i.outBuf[f]
		t++
	}
	i.outEnd = t
}

func (i *IconvReadCloser) moveOrgBuffer(end int) {
	if end == 0 {
		return
	}
	t := 0
	for f := end; f < i.inEnd; f++ {
		i.inBuf[t] = i.inBuf[f]
		t++
	}
	i.inEnd = t
}

type IconvWriteCloser struct {
	iconv  Iconv
	w      io.Writer
	outBuf []byte
}

func NewIconvWriteCloser(w io.Writer, tocode, fromcode string) (*IconvWriteCloser, error) {
	return NewIconvWriteCloserBufferSize(w, 1024, tocode, fromcode)
}

func NewIconvWriteCloserBufferSize(w io.Writer, bufSize int, tocode, fromcode string) (*IconvWriteCloser, error) {
	iconv, err := NewIconv(tocode, fromcode)
	if err != nil {
		return nil, err
	}
	return &IconvWriteCloser{
		iconv:  iconv,
		w:      w,
		outBuf: make([]byte, bufSize, bufSize),
	}, nil
}

func (i *IconvWriteCloser) Close() error {
	return i.iconv.Close()
}

func (i *IconvWriteCloser) Write(p []byte) (int, error) {
	for l, mp := 0, len(p); l < mp; {
		in := p[l:]
		out := i.outBuf[:]
		inlen, outlen, err := i.iconv.Conv(in, out)
		if err != nil && outlen == 0 {
			if l == 0 {
				fmt.Println("err with 0")
				return l, err
			} else {
				return l, nil
			}
		}
		out = i.outBuf[:outlen]
		for len(out) > 0 {
			n, err := i.w.Write(out)
			if err != nil {
				return l, err
			}
			out = out[n:]
		}
		l += inlen
	}
	return len(p), nil
}
