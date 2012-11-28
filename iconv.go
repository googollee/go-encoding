package encodingex

// #cgo LDFLAGS: -liconv
// #include <iconv.h>
import "C"

import (
	"unsafe"
)

type Iconver struct {
	p C.iconv_t
}

func NewIconver(tocode, fromcode string) (i Iconver, err error) {
	i.p, err = C.iconv_open(C.CString(tocode), C.CString(fromcode))
	return
}

func (i Iconver) Conv(in, out []byte) (inlen int, outlen int, err error) {
	insize, outsize := C.size_t(len(in)), C.size_t(len(out))
	inptr, outptr := &in[0], &out[0]
	_, err = C.iconv(i.p,
		(**C.char)(unsafe.Pointer(&inptr)), &insize,
		(**C.char)(unsafe.Pointer(&outptr)), &outsize)
	inlen, outlen = len(in)-int(insize), len(out)-int(outsize)
	return
}

func (i Iconver) Close() (err error) {
	_, err = C.iconv_close(i.p)
	return
}
