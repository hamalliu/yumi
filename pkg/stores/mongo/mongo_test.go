package mongo

import "testing"

func TestConnect(t *testing.T) {
	err := Connect()
	if err != nil {
		t.Log(err)
	}
}
