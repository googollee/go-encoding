package encodingex

import (
	"io"
	"fmt"
)

func unHex(c byte) (byte, bool) {
	if byte('0') <= c && c <= byte('9') {
		return c - byte('0'), true
	}
	if byte('a') <= c && c <= byte('f') {
		return c - byte('a') + 10, true
	}
	if byte('A') <= c && c <= byte('F') {
		return c - byte('A') + 10, true
	}
	return 0, false
}

func hex(p byte) (ret [2]byte) {
	ret[0] = p / 16
	if ret[0] >= 10 {
		ret[0] = ret[0] - 10 + byte('A')
	} else {
		ret[0] = ret[0] + byte('0')
	}
	ret[1] = p % 16
	if ret[1] >= 10 {
		ret[1] = ret[1] - 10 + byte('A')
	} else {
		ret[1] = ret[1] + byte('0')
	}
	return
}

type decodeStatus int

const (
	decodeQuoted decodeStatus = iota
	decodeFirst
	decodeReturn
	decodeNormal
)

type decoder struct {
	r      io.Reader
	status decodeStatus
	temp   byte
}

func NewQuotedPrintableDecoder(r io.Reader) io.Reader {
	return &decoder{
		r:      r,
		status: decodeNormal,
		temp:   0,
	}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	readed := make([]byte, len(p), cap(p))
	readedn, err := d.r.Read(readed)
	if err != nil {
		return
	}
	n = 0
	for _, b := range readed[:readedn] {
		switch d.status {
		case decodeNormal:
			if b == byte('=') {
				d.status = decodeQuoted
			} else {
				p[n] = b
				n++
			}
		case decodeQuoted:
			switch b {
			case byte('\n'):
				d.status = decodeNormal
			case byte('\r'):
				d.status = decodeReturn
			default:
				h, ok := unHex(b)
				if !ok {
					err = fmt.Errorf("can't convert %c(%d) to hex", rune(b), b)
					return
				}
				d.temp = h * 16
				d.status = decodeFirst
			}
		case decodeReturn:
			if b != byte('\n') {
				p[n] = d.temp
				n++
			}
			d.status = decodeNormal
		case decodeFirst:
			h, ok := unHex(b)
			if !ok {
				err = fmt.Errorf("can't convert %c(%d) to hex", rune(b), b)
				return
			}
			d.temp += h
			p[n] = d.temp
			n++
			d.status = decodeNormal
		}
	}
	return
}

type encodeStatus int

const (
	encodeNormal encodeStatus = iota
	encodeSpace
	encodeSpaceReturn
	encodeReturn
)

const maxBuf = 1024

type encoder struct {
	w             io.Writer
	status        encodeStatus
	maxLineLength int
	lineLength    int
	nbuf          int
	buf           [maxBuf]byte
	last          byte
}

func NewQuotedPrintableEncoder(w io.Writer, maxLength int) io.WriteCloser {
	if maxLength < 3 {
		maxLength = 3
	}
	if maxLength > 76 {
		maxLength = 76
	}
	return &encoder{
		w:             w,
		status:        encodeNormal,
		maxLineLength: maxLength,
		lineLength:    0,
		nbuf:          0,
	}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	e.nbuf = 0
	var b byte
	for n, b = range p {
		switch e.status {
		case encodeNormal:
			switch b {
			case byte(' '):
				fallthrough
			case byte('\t'):
				e.status = encodeSpace
				e.last = b
			case byte('\r'):
				e.status = encodeReturn
			case byte('\n'):
				if err = e.push(byte('\r')); err != nil {
					return
				}
				if err = e.push(byte('\n')); err != nil {
					return
				}
			case byte('='):
				if err = e.pushQuoted(byte('=')); err != nil {
					return
				}
			default:
				if err = e.pushCheck(b); err != nil {
					return
				}
			}
		case encodeSpace:
			switch b {
			case byte('\r'):
				e.status = encodeSpaceReturn
			case byte('\n'):
				e.status = encodeNormal
				if err = e.pushQuoted(e.last); err != nil {
					return
				}
				if err = e.push('\r'); err != nil {
					return
				}
				if err = e.push('\n'); err != nil {
					return
				}
			default:
				e.status = encodeNormal
				if err = e.push(e.last); err != nil {
					return
				}
				if err = e.pushCheck(b); err != nil {
					return
				}
			}
		case encodeSpaceReturn:
			if b == byte('\n') {
				if err = e.pushQuoted(e.last); err != nil {
					return
				}
				if err = e.push('\r'); err != nil {
					return
				}
				if err = e.push('\n'); err != nil {
					return
				}
			} else {
				if err = e.push(e.last); err != nil {
					return
				}
				if err = e.pushQuoted(byte('\r')); err != nil {
					return
				}
				if err = e.pushCheck(b); err != nil {
					return
				}
			}
		case encodeReturn:
			if b == byte('\n') {
				if err = e.push(byte('\r')); err != nil {
					return
				}
				if err = e.push(byte('\n')); err != nil {
					return
				}
			} else {
				if err = e.pushQuoted(byte('\r')); err != nil {
					return
				}
				if err = e.pushCheck(b); err != nil {
					return
				}
			}
		}
	}
	_, err = e.w.Write(e.buf[:e.nbuf])
	return
}

func (e *encoder) Close() error {
	return nil
}

func (e *encoder) pushCheck(p byte) error {
	if 33 <= p && p <= 126 {
		return e.push(p)
	}
	return e.pushQuoted(p)
}

func (e *encoder) push(p byte) error {
	if p != byte('\r') && p != byte('\n') {
		if (e.lineLength + 1) >= e.maxLineLength {
			e.buf[e.nbuf] = byte('=')
			e.nbuf++
			e.buf[e.nbuf] = byte('\r')
			e.nbuf++
			e.buf[e.nbuf] = byte('\n')
			e.nbuf++
			e.lineLength = 0
		}
	} else {
		e.lineLength = 0
	}
	e.buf[e.nbuf] = p
	e.nbuf++
	e.lineLength++
	return e.checkAndSendBuffer()
}

func (e *encoder) pushQuoted(p byte) error {
	if (e.lineLength + 3) >= e.maxLineLength {
		e.push(byte('='))
		e.push(byte('\r'))
		e.push(byte('\n'))
	}
	if err := e.push(byte('=')); err != nil {
		return err
	}
	for _, c := range hex(p) {
		if err := e.push(byte(c)); err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) checkAndSendBuffer() error {
	if e.nbuf >= maxBuf {
		_, err := e.w.Write(e.buf[:])
		if err != nil {
			return err
		}
		e.nbuf = 0
	}
	return nil
}
