package utils

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func UTF2GBK(str string) string {
	t := transform.NewReader(bytes.NewReader([]byte(str)), simplifiedchinese.GBK.NewEncoder())
	data, _ := ioutil.ReadAll(t)
	return string(data)
}
