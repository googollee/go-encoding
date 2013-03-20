package encoding

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestIgnoreReader(t *testing.T) {
	str := "\r\n1231\tjkl\rrewq\nerq"
	expect := "1231\tjklrewqerq"
	reader := NewIgnoreReader(bytes.NewBufferString(str), []byte("\r\n"))
	got, _ := ioutil.ReadAll(reader)
	if string(got) != expect {
		t.Errorf("Expect: %s, got: %s", expect, string(got))
	}
}

func TestIgnoreWrite(t *testing.T) {
	input := []byte("\r\n1231\tjkl\rrewq\nerq")
	expect := "1231\tjklrewqerq"
	buf := bytes.NewBuffer(nil)
	writer := NewIgnoreWriter(buf, []byte("\r\n"))
	n, _ := writer.Write(input)
	if n != len(input) {
		t.Errorf("Should write %d bytes, actual write %d bytes", len(input), n)
	}
	if buf.String() != expect {
		t.Errorf("Expect: %s, got: %s", expect, buf.String())
	}
}
