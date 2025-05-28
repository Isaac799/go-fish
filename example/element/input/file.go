package element

import (
	"fmt"
	"mime/multipart"
	"strings"
)

func InputJoinComma(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, ",")
}

type InputFileValue struct {
	file   multipart.File
	header *multipart.FileHeader
}

type InputFile struct {
	ID       string
	Label    string
	Key      string
	Value    *InputFileValue
	Accept   []string
	Multiple bool
	Disabled bool
	Required bool
}

func NewInputFile(name string) InputFile {
	return InputFile{
		ID:       fmt.Sprintf("id-%s", name),
		Label:    name,
		Key:      name,
		Value:    nil,
		Accept:   []string{},
		Multiple: true,
		Disabled: false,
		Required: false,
	}
}
