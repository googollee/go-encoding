package encodingex

import (
	"testing"
)

func TestIconv(t *testing.T) {
	str := "你好世界"
	temp1 := make([]byte, 20, 20)
	temp2 := make([]byte, 20, 20)

	from, err := NewIconv("gbk", "utf-8")
	t.Logf("%s", err)
	defer from.Close()
	to, err := NewIconv("utf-8", "gbk")
	t.Logf("%s", err)
	defer to.Close()

	in, out, err := from.Conv([]byte(str), temp1[:])
	t.Logf("%s", err)
	if expect := 12; in != expect {
		t.Errorf("expect: %d, got: %d", expect, in)
	}
	if expect := 8; out != expect {
		t.Errorf("expect: %d, got: %d", expect, out)
	}
	t.Logf("%s", string(temp1[:out]))

	in, out, err = to.Conv(temp1[:out], temp2[:])
	t.Logf("%s", err)
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
