package encodingex

import (
	"testing"
	"io/ioutil"
	"bytes"
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
