package encoding

import (
	"testing"
)

func TestEncodedWordDeocder(t *testing.T) {
	{
		str := "=?iso-8859-1?Q?=A1Hola,_se=F1or!?="
		expect := "¡Hola, señor!"
		got, _ := DecodeEncodedWord(str)
		if got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}

	{
		str := "=?GB2312?B?W1RMXSBwb25nYmFAZ29vZ2xlZ3JvdXBzLmNvbSC1xNWq0qogLSChsDEguPbW9w==?= \n=?GB2312?B?zOKhsdPQIDEguPbM+9fT?="
		expect := "[TL] pongba@googlegroups.com 的摘要 - “1 个主题”有 1 个帖子"
		got, _ := DecodeEncodedWord(str)
		if got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}

	{
		strs := []string{"=?GB2312?B?W1RMXSBwb25nYmFAZ29vZ2xlZ3JvdXBzLmNvbSC1xNWq0qogLSChsDEguPbW9w==?=", "=?GB2312?B?zOKhsdPQIDEguPbM+9fT?="}
		expect := "[TL] pongba@googlegroups.com 的摘要 - “1 个主题”有 1 个帖子"
		got, _ := DecodeEncodedWordArray(strs)
		if got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}
}
