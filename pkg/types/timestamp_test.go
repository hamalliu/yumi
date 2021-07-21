package types

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTimestamp(t *testing.T) {
	type AA struct {
		A Timestamp `json:"a"`
	}

	type Data struct {
		A    string      `json:"a"`
		Data interface{} `json:"data"`
	}

	a := AA{A: 1626415470}
	d := Data{Data: &a, A: "xxxx"}

	err := json.NewEncoder(os.Stdout).Encode(&d)
	if err != nil {
		t.Error(err)
	}
}
