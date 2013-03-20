package encoding

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
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

// Decode an encoded-word string. The string can be a set of encoded-word splited by CRLF SPACE.
func DecodeEncodedWord(str string) (ret string, err error) {
	replacer := strings.NewReplacer("\r\n", " ", "\n", " ")
	str = replacer.Replace(str)
	return DecodeEncodedWordArray(strings.Split(str, " "))
}

// Decode an array of encoded-word string
func DecodeEncodedWordArray(strs []string) (string, error) {
	if len(strs) == 0 {
		return "", nil
	}
	splits, err := splitString(strs[0])
	if err != nil {
		return "", fmt.Errorf("%s: line 0", err)
	}
	charset := strings.ToUpper(splits[1])
	codec := strings.ToUpper(splits[2])
	bufs := make([]io.Reader, 0)
	for i, s := range strs {
		if s == "" {
			continue
		}
		var splits []string
		splits, err = splitString(s)
		if err != nil {
			return "", fmt.Errorf("%s: line %d", err, i)
		}
		if strings.ToUpper(splits[1]) != charset {
			return "", fmt.Errorf("line %d charset invalid: %s", i, splits[1])
		}
		if strings.ToUpper(splits[2]) != codec {
			return "", fmt.Errorf("line %d codec invalid: %s", i, splits[2])
		}
		bufs = append(bufs, bytes.NewBufferString(splits[3]))
	}
	var reader io.Reader = io.MultiReader(bufs...)
	switch codec {
	case "Q":
		reader = NewQEncodingDecoder(reader)
	case "B":
		reader = base64.NewDecoder(base64.StdEncoding, reader)
	default:
		return "", fmt.Errorf("invalid codec: %s", codec)
	}
	if charset != "UTF-8" {
		r, err := NewIconvReadCloser(reader, "UTF-8", charset)
		if err != nil {
			return "", err
		}
		defer r.Close()
		reader = r
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
