package encodingex

import (
	"io"
	"io/ioutil"
	"strings"
	"encoding/base64"
	"errors"
	"fmt"
	"bytes"
)

func splitString(str string) ([]string, error) {
	splits := strings.Split(str, "?")
	if len(splits) != 5 {
		return nil, errors.New("invalid format")
	}
	if splits[0] != "=" && splits[4] != "=" {
		return nil, errors.New("invalid format")
	}
	return splits, nil
}

func DecodeEncodedWord(str string) (ret, charset string, err error) {
	var splits []string
	splits, err = splitString(str)
	if err != nil {
		return
	}
	charset = strings.ToUpper(splits[1])
	codec := strings.ToUpper(splits[2])
	buf := bytes.NewBufferString(splits[3])
	var reader io.Reader
	switch codec {
	case "Q":
		reader = NewQEncodingDecoder(buf)
	case "B":
		reader = base64.NewDecoder(base64.StdEncoding, buf)
	default:
		err = fmt.Errorf("invalid codec: %s", codec)
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	ret = string(data)
	return
}

func DecodeEncodedWordArray(strs []string) (ret, charset string, err error) {
	if len(strs) == 0 {
		return
	}
	var splits []string
	splits, err = splitString(strs[0])
	if err != nil {
		err = fmt.Errorf("%s: line 0", err)
		return
	}
	charset = strings.ToUpper(splits[1])
	codec := strings.ToUpper(splits[2])
	bufs := make([]io.Reader, len(strs), len(strs))
	for i, s := range strs {
		var splits []string
		splits, err = splitString(s)
		if err != nil {
			err = fmt.Errorf("%s: line %d", err, i)
			return
		}
		bufs[i] = bytes.NewBufferString(splits[3])
	}
	buf := io.MultiReader(bufs...)
	var reader io.Reader
	switch codec {
	case "Q":
		reader = NewQEncodingDecoder(buf)
	case "B":
		reader = base64.NewDecoder(base64.StdEncoding, buf)
	default:
		err = fmt.Errorf("invalid codec: %s", codec)
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	ret = string(data)
	return
}
