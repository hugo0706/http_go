package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"
var fieldNameFormat = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_|~` + "`" + `]+$`)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}
	fieldLine := string(data[:idx])
	parts := strings.SplitN(fieldLine, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("Invalid header")
	}
	fieldNameInRunes := []rune(parts[0])
	if fieldNameInRunes[len(fieldNameInRunes)-1] == ' ' {
		return 0, false, errors.New("Invalid header")
	}

	fieldName := strings.Trim(parts[0], " ")
	if m := fieldNameFormat.MatchString(fieldName); !m {
		return 0, false, errors.New("Invalid character found on header")
	}

	fieldValue := strings.Trim(parts[1], " ")
	fieldName = strings.ToLower(fieldName)

	h.Set(fieldName, fieldValue)

	return len(fieldLine) + len(crlf), false, nil
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v,value}, ", ")
	}
	h[key] = value
}