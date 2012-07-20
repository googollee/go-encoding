package encodingex

import (
	"io"
)

type IconvReadCloser struct {
	r      io.Reader
	iconv  Iconv
	maxBuf int
	orgBuf []byte
	orgEnd int
	outBuf []byte
	outEnd int
}

func NewIconvReadCloser(r io.Reader, tocode, fromcode string) (*IconvReadCloser, error) {
	return NewIconvReadCloserBufferSize(r, maxBuf, tocode, fromcode)
}

func NewIconvReadCloserBufferSize(r io.Reader, maxBuf int, tocode, fromcode string) (*IconvReadCloser, error) {
	iconv, err := NewIconv(tocode, fromcode)
	if err != nil {
		return nil, err
	}
	return &IconvReadCloser{
		r:      r,
		iconv:  iconv,
		maxBuf: maxBuf,
		orgBuf: make([]byte, maxBuf, maxBuf),
		orgEnd: 0,
		outBuf: make([]byte, maxBuf, maxBuf),
		outEnd: 0,
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
	if i.outEnd == maxBuf {
		return nil
	}
	if i.orgEnd < maxBuf {
		n, err := i.r.Read(i.orgBuf[i.orgEnd:])
		if err != nil {
			return err
		}
		i.orgEnd += n
	}
	inbuf := i.orgBuf[:i.orgEnd]
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
	for f := end; f < i.orgEnd; f++ {
		i.orgBuf[t] = i.orgBuf[f]
		t++
	}
	i.orgEnd = t
}
