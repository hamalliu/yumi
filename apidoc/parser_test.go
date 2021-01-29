package apidoc

import "testing"

func TestParseHandle(t *testing.T) {
	docs, err := ParseHandleDocs("../api/media")
	if err != nil {
		panic(err)
	}
	t.Log(docs)
}

func TestParse(t *testing.T) {
	t.Log(Parse("../api", "/"))
}
