package encodingex

import (
	"io"
)

type iconver struct {
	iconv  Iconv
	maxBuf int
	inBuf  []byte
	inEnd  int
	outBuf []byte
	outEnd int
}

func newIconver(maxBuf int, tocode, fromcode string) (i iconver, err error) {
	iconv, e := NewIconv(tocode, fromcode)
	if e != nil {
		err = e
		return
	}
	i.iconv = iconv
	i.maxBuf = maxBuf
	i.inBuf = make([]byte, maxBuf, maxBuf)
	i.inEnd = 0
	i.outBuf = make([]byte, maxBuf, maxBuf)
	i.outEnd = 0
	return
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

type IconvReadCloser struct {
	iconver
	r io.Reader
}

func NewIconvReadCloser(r io.Reader, tocode, fromcode string) (*IconvReadCloser, error) {
	return NewIconvReadCloserBufferSize(r, maxBuf, tocode, fromcode)
}

func NewIconvReadCloserBufferSize(r io.Reader, maxBuf int, tocode, fromcode string) (*IconvReadCloser, error) {
	iconver, err := newIconver(maxBuf, tocode, fromcode)
	if err != nil {
		return nil, err
	}
	return &IconvReadCloser{iconver, r}, nil
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
	if i.outEnd == i.maxBuf {
		return nil
	}
	if i.inEnd < i.maxBuf {
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

type IconvWriteCloser struct {
	iconver
	w io.WriteCloser
}

func NewIconvWriteCloser(w io.WriteCloser, tocode, fromcode string) (*IconvWriteCloser, error) {
	return NewIconvWriteCloserBufferSize(w, maxBuf, tocode, fromcode)
}

func NewIconvWriteCloserBufferSize(w io.WriteCloser, maxBuf int, tocode, fromcode string) (*IconvWriteCloser, error) {
	iconver, err := newIconver(maxBuf, tocode, fromcode)
	if err != nil {
		return nil, err
	}
	return &IconvWriteCloser{iconver, w}, nil
}

func (i *IconvWriteCloser) Close() error {
	err := i.w.Close()
	if err != nil {
		return err
	}
	return i.iconv.Close()
}
