package encodingex

import (
	"fmt"
	"io"
)

type qeDecodeStatus int

const (
	qeDecodeQuoted qeDecodeStatus = iota
	qeDecodeFirst
	qeDecodeReturn
	qeDecodeNormal
)

type qeDecoder struct {
	r      io.Reader
	status qeDecodeStatus
	temp   byte
}

// Create a Reader to decode q-encoding from r.
func NewQEncodingDecoder(r io.Reader) io.Reader {
	return &qeDecoder{
		r:      r,
		status: qeDecodeNormal,
		temp:   0,
	}
}

func (d *qeDecoder) Read(p []byte) (n int, err error) {
	readed := make([]byte, len(p), cap(p))
	readedn, err := d.r.Read(readed)
	if err != nil {
		return
	}
	n = 0
	for _, b := range readed[:readedn] {
		switch d.status {
		case qeDecodeNormal:
			switch b {
			case byte('='):
				d.status = qeDecodeQuoted
			case byte('_'):
				b = byte(' ')
				fallthrough
			default:
				p[n] = b
				n++
			}
		case qeDecodeQuoted:
			switch b {
			case byte('\n'):
				d.status = qeDecodeNormal
			case byte('\r'):
				d.status = qeDecodeReturn
			default:
				h, ok := unHex(b)
				if !ok {
					err = fmt.Errorf("can't convert %c(%d) to hex", rune(b), b)
					return
				}
				d.temp = h * 16
				d.status = qeDecodeFirst
			}
		case qeDecodeReturn:
			if b != byte('\n') {
				p[n] = d.temp
				n++
			}
			d.status = qeDecodeNormal
		case qeDecodeFirst:
			h, ok := unHex(b)
			if !ok {
				err = fmt.Errorf("can't convert %c(%d) to hex", rune(b), b)
				return
			}
			d.temp += h
			p[n] = d.temp
			n++
			d.status = qeDecodeNormal
		}
	}
	return
}
