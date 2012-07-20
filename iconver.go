package encodingex

import (
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
	i.arrangeOutBuffer(m)
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
	inbuf, outbuf := i.inBuf[:i.inEnd], i.outBuf[i.outEnd:]
	inlen, outlen, err := i.iconv.Conv(inbuf, outbuf)
	i.arrangeInBuffer(inlen)
	i.outEnd += outlen
	return err
}

func (i *IconvReadCloser) arrangeOutBuffer(index int) {
	if index == 0 {
		return
	}
	end := i.outEnd
	i.outEnd = 0
	for ; index < end; index++ {
		i.outBuf[i.outEnd] = i.outBuf[index]
		i.outEnd++
	}
}

func (i *IconvReadCloser) arrangeInBuffer(index int) {
	if index == 0 {
		return
	}
	end := i.inEnd
	i.inEnd = 0
	for ; index < end; index++ {
		i.inBuf[i.inEnd] = i.inBuf[index]
		i.inEnd++
	}
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
	for index := 0; index < len(p); {
		in, out := p[index:], i.outBuf[:]
		inlen, outlen, err := i.iconv.Conv(in, out)
		if err != nil && outlen == 0 {
			if index == 0 {
				return index, err
			} else {
				return index, nil
			}
		}

		for out := i.outBuf[:outlen]; len(out) > 0; {
			n, err := i.w.Write(out)
			if err != nil {
				return index, err
			}
			out = out[n:]
		}
		index += inlen
	}
	return len(p), nil
}
