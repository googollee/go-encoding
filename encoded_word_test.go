package encodingex

import (
	"testing"
)

func TestEncodedWordDeocder(t *testing.T) {
	{
		str := "=?iso-8859-1?Q?=A1Hola,_se=F1or!?="
		expect := "¡Hola, señor!"
		got, charset, _ := DecodeEncodedWord(str)
		got, _ = Conv(got, "UTF-8", charset)
		if got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}

	{
		strs := []string{"=?GB2312?B?W1RMXSBwb25nYmFAZ29vZ2xlZ3JvdXBzLmNvbSC1xNWq0qogLSChsDEguPbW9w==?=", "=?GB2312?B?zOKhsdPQIDEguPbM+9fT?="}
		expect := "[TL] pongba@googlegroups.com 的摘要 - “1 个主题”有 1 个帖子"
		got, charset, _ := DecodeEncodedWordArray(strs)
		got, _ = Conv(got, "UTF-8", charset)
		if got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}
}
