package encodingex

// #cgo freebsd LDFLAGS: -L/usr/local/lib -liconv
// #cgo freebsd CFLAGS: -I/usr/local/include
// #cgo darwin LDFLAGS: -liconv
// #include <iconv.h>
import "C"

import (
	"unsafe"
)

type Iconv struct {
	p C.iconv_t
}

// Create a Iconv instance, convert codec from fromcode to tocode.
func NewIconv(tocode, fromcode string) (i Iconv, err error) {
	i.p, err = C.iconv_open(C.CString(tocode), C.CString(fromcode))
	return
}

// Do convert from in to out.
func (i Iconv) Conv(in, out []byte) (inlen int, outlen int, err error) {
	insize, outsize := C.size_t(len(in)), C.size_t(len(out))
	inptr, outptr := &in[0], &out[0]
	_, err = C.iconv(i.p,
		(**C.char)(unsafe.Pointer(&inptr)), &insize,
		(**C.char)(unsafe.Pointer(&outptr)), &outsize)
	inlen, outlen = len(in)-int(insize), len(out)-int(outsize)
	return
}

// Close a Iconv.
func (i Iconv) Close() (err error) {
	_, err = C.iconv_close(i.p)
	return
}
