package encoding

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestQPrintableDecode(t *testing.T) {
	{
		str := "If_you_believe_that_truth=3Dbeauty,_then=5Fsurely_=\nmathematics_is_the_most_beautiful_branch_of_phil=\nosophy.=20\nabc=3F"
		expect := "If you believe that truth=beauty, then_surely mathematics is the most beautiful branch of philosophy. \nabc?"
		reader := NewQEncodingDecoder(bytes.NewBufferString(str))
		got, _ := ioutil.ReadAll(reader)
		if string(got) != expect {
			t.Errorf("Decode\n\texpect: %s\n\t   got: %s", expect, string(got))
		}
	}

	{
		str := "If_you_believe_that_truth=3Dbeauty,_then=5Fsurely_=\r\nmathematics_is_the_most_beautiful_branch_of_phil=\r\nosophy.=20\r\nabc=3F"
		expect := "If you believe that truth=beauty, then_surely mathematics is the most beautiful branch of philosophy. \r\nabc?"
		reader := NewQEncodingDecoder(bytes.NewBufferString(str))
		got, err := ioutil.ReadAll(reader)
		t.Log(err)
		if string(got) != expect {
			t.Errorf("Decode\n\texpect: %s\n\t   got: %s", expect, string(got))
		}
	}
}
