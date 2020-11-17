package html2docx

import (
	"strings"
)

// Rune ...
type Rune struct {
	// Text 文本
	Text    string
	HTMLTag string
	start   int
	end     int
}

// Tag ...
type Tag struct {
}

var (
	prefixTags = []string{
		"<h1>",
	}

	suffixTags = []string{
		"</h1>",
	}
)

// var str := "《手动《手动《手动阀手动阀》阀手《发生的人格》动阀》阀手《粉色粉色我》动阀》"

// TODO: 此方法效率较低，需优化
func getfirstTag(text string) (tagIndex, taglen int, tagData string) {
	tagIndex = -1
	l := len(prefixTags) + len(suffixTags)
	tags := make([]string, l, l)
	tags = append(tags, prefixTags...)
	tags = append(tags, suffixTags...)
	for _, tag := range tags {
		index := strings.Index(text, tag)
		if index != -1 && (tagIndex == -1 || tagIndex > index) {
			tagIndex = index
			taglen = len(tag)
			tagData = tag
		}
	}

	return
}

//
func getLeftTags(rs []Rune, text string) []Rune {
	// for i := range rs {
	// 	rs[i].Text
	// }

	return nil
}
