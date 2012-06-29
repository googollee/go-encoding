package goqp

import (
	"io"
	"fmt"
)

func UnHex(c byte) (byte, bool) {
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

type decodeStatus int

const (
	decodeQuoted decodeStatus = iota
	decodeFirst
	decodeReturn
	decodeNormal
)

type decoder struct {
	r io.Reader
	status decodeStatus
	temp byte
}

func NewDecoder(r io.Reader) io.Reader {
	return &decoder{
		r: r,
		status: decodeNormal,
		temp: 0,
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
				h, ok := UnHex(b)
				if !ok {
					err = fmt.Errorf("can't convert %c(%d) to hex", rune(b), b)
					return
				}
				d.temp = h * 16
				d.status = decodeFirst
			}
		case decodeReturn:
			if b == byte('\n') {
				d.status = decodeNormal
			} else {
				p[n] = d.temp
				n++
			}
		case decodeFirst:
			h, ok := UnHex(b)
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