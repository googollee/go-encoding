package encodingex

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestIconver(t *testing.T) {
	str := "你好世界"
	temp1 := make([]byte, 20, 20)
	temp2 := make([]byte, 20, 20)

	from, err := NewIconver("gbk", "utf-8")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer from.Close()
	to, err := NewIconver("utf-8", "gbk")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer to.Close()

	in, out, err := from.Conv([]byte(str), temp1[:])
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if expect := 12; in != expect {
		t.Errorf("expect: %d, got: %d", expect, in)
	}
	if expect := 8; out != expect {
		t.Errorf("expect: %d, got: %d", expect, out)
	}
	t.Logf("%s", string(temp1[:out]))

	in, out, err = to.Conv(temp1[:out], temp2[:])
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if expect := 8; in != expect {
		t.Errorf("expect: %d, got: %d", expect, in)
	}
	if expect := 12; out != expect {
		t.Errorf("expect: %d, got: %d", expect, out)
	}
	if got := string(temp2[:out]); got != str {
		t.Errorf("expect: %s, got: %s", str, got)
	}
	t.Logf("%s", string(temp2[:out]))
}

func TestIconvPart(t *testing.T) {
	str := "你好世界"
	temp1 := make([]byte, 20, 20)
	temp2 := make([]byte, 20, 20)

	from, err := NewIconver("gbk", "utf-8")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer from.Close()
	to, err := NewIconver("utf-8", "gbk")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer to.Close()

	in, out, err := from.Conv([]byte(str)[:10], temp1[:])
	if err == nil {
		t.Errorf("expect err not nil")
	}
	if expect := 9; in != expect {
		t.Errorf("expect: %d, got: %d", expect, in)
	}
	if expect := 6; out != expect {
		t.Errorf("expect: %d, got: %d", expect, out)
	}
	t.Logf("%s", string(temp1[:out]))

	in, out, err = to.Conv(temp1[:out], temp2[:])
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if expect := 6; in != expect {
		t.Errorf("expect: %d, got: %d", expect, in)
	}
	if expect := 9; out != expect {
		t.Errorf("expect: %d, got: %d", expect, out)
	}
	if got, expect := string(temp2[:out]), "你好世"; got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
	t.Logf("%s", string(temp2[:out]))
}

func TestIconvReadCloser(t *testing.T) {
	str := ""
	for i := 0; i < 200; i++ {
		str = str + "你好世界"
	}
	buf := bytes.NewBufferString(str)

	from, err := NewIconvReadCloser(buf, "gbk", "utf-8")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer from.Close()
	to, err := NewIconvReadCloserBufferSize(from, 233, "utf-8", "gbk")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer to.Close()

	temp, err := ioutil.ReadAll(to)
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if expect, got := 2400, len(temp); got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	if expect, got := str, string(temp); got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
	_, err = to.Read(temp)
	if err != io.EOF {
		fmt.Println("expect eof, got: %s", err)
	}
}

func TestIconvWriteCloser(t *testing.T) {
	str := ""
	for i := 0; i < 200; i++ {
		str = str + "你好世界"
	}
	buf := bytes.NewBuffer(nil)

	from, err := NewIconvWriteCloser(buf, "gbk", "utf-8")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer from.Close()
	to, err := NewIconvWriteCloser(from, "utf-8", "gbk")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	defer to.Close()

	n, err := to.Write([]byte(str))
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if expect, got := 2400, n; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	if expect, got := str, buf.String(); got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
}

func TestIconv(t *testing.T) {
	str := "你好世界"
	to, err := Iconv(str, "gbk", "utf-8")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if str == to {
		t.Errorf("expect: %s not same as %s", to, str)
	}

	from, err := Iconv(to, "utf-8", "gbk")
	if err != nil {
		t.Errorf("expect err nil, got: %s", err)
	}
	if str != from {
		t.Errorf("expect: %s same as %s", from, str)
	}
}
