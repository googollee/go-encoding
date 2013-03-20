package encoding

import (
	"bytes"
	"io/ioutil"
	"testing"
)

type unHexData struct {
	in  byte
	out byte
	ok  bool
}

func TestUnHex(t *testing.T) {
	datas := [...]unHexData{
		{byte('0'), 0, true},
		{byte('6'), 6, true},
		{byte('a'), 10, true},
		{byte('d'), 13, true},
		{byte('A'), 10, true},
		{byte('D'), 13, true},
		{byte('y'), 0, false},
		{byte('Y'), 0, false},
		{byte('$'), 0, false},
	}
	for _, d := range datas {
		got, ok := unHex(d.in)
		if ok != d.ok {
			t.Errorf("convert %d(%c) should be %v, got: %v", d.in, rune(d.in), d.ok, ok)
		}
		if ok && got != d.out {
			t.Errorf("convert %d(%c) should get %d(%c), got: %d(%c)",
				d.in, rune(d.in), d.out, rune(d.out), got, rune(got))
		}
	}
}

type hexData struct {
	in  byte
	out [2]byte
}

func TestHex(t *testing.T) {
	datas := [...]hexData{
		{byte('='), [2]byte{byte('3'), byte('D')}},
		{byte(' '), [2]byte{byte('2'), byte('0')}},
		{byte('\r'), [2]byte{byte('0'), byte('D')}},
	}
	for _, d := range datas {
		got := hex(d.in)
		if got != d.out {
			t.Errorf("convert %d(%c) should be %v, got: %v", d.in, rune(d.in), d.out, got)
		}
	}
}

func TestQuotedPrintableDecode(t *testing.T) {
	{
		str := "If you believe that truth=3Dbeauty, then surely =\nmathematics is the most beautiful branch of phil=\nosophy.=20\nabc"
		expect := "If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy. \nabc"
		reader := NewQuotedPrintableDecoder(bytes.NewBufferString(str))
		got, _ := ioutil.ReadAll(reader)
		if string(got) != expect {
			t.Errorf("Decode\n\texpect: %s\n\t   got: %s", expect, string(got))
		}
	}

	{
		str := "If you believe that truth=3Dbeauty, then surely =\r\nmathematics is the most beautiful branch of phil=\r\nosophy.=20\r\nabc"
		expect := "If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy. \r\nabc"
		reader := NewQuotedPrintableDecoder(bytes.NewBufferString(str))
		got, err := ioutil.ReadAll(reader)
		t.Log(err)
		if string(got) != expect {
			t.Errorf("Decode\n\texpect: %s\n\t   got: %s", expect, string(got))
		}
	}
}

func TestQuotedPrintableEncode(t *testing.T) {
	{
		input := []byte("If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy. \nabc")
		expect := "If you believe that truth=3Dbeauty, then surely =\r\nmathematics is the most beautiful branch of phil=\r\nosophy.=20\r\nabc"
		buf := bytes.NewBuffer(nil)
		writer := NewQuotedPrintableEncoder(buf, 49)
		_, _ = writer.Write(input)
		if buf.String() != expect {
			t.Errorf("Encode\n\texpect: %s\n\t   got: %s", expect, buf.String())
		}
	}
}
