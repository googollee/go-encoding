package encodingex

// #cgo LDFLAGS: -liconv
// #include <iconv.h>
import "C"

import (
	"unsafe"
)

type Iconv struct {
	p C.iconv_t
}

func NewIconv(tocode, fromcode string) (Iconv, error) {
	i, err := C.iconv_open(C.CString(tocode), C.CString(fromcode))
	if err != nil {
		return Iconv{}, err
	}
	return Iconv{i}, nil
}

func (i Iconv) Conv(in, out []byte) (inlen int, outlen int, err error) {
	insize := C.size_t(len(in))
	outsize := C.size_t(len(out))
	inptr := &in[0]
	outptr := &out[0]
	_, err = C.iconv(i.p,
		(**C.char)(unsafe.Pointer(&inptr)), &insize,
		(**C.char)(unsafe.Pointer(&outptr)), &outsize)
	inlen = len(in) - int(insize)
	outlen = len(out) - int(outsize)
	return
}

func (i Iconv) Close() error {
	_, err := C.iconv_close(i.p)
	if err != nil {
		return err
	}
	return nil
}
