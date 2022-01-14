package types

import (
	"encoding/json"
	"testing"
)

func TestIP(t *testing.T) {
	str := `{
		"a": "10.34.1.1",
		"b": "10.24.1.1"
	}`

	ret := struct {
		A IP `json:"a"`
		B IP `json:"b"`
	}{}

	if err := json.Unmarshal([]byte(str), &ret); err != nil {
		t.Error(err)
		return
	}
	retBs, err := ret.A.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret, string(retBs))
}
