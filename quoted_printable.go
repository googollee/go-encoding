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
			if b == byte('\n') {
				d.status = decodeNormal
			} else {
				h, ok := UnHex(b)
				if !ok {
					err = fmt.Errorf("can't convert %c to hex", rune(b))
					return
				}
				d.temp = h * 16
				d.status = decodeFirst
			}
		case decodeFirst:
			h, ok := UnHex(b)
			if !ok {
				err = fmt.Errorf("can't convert %c to hex", rune(b))
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