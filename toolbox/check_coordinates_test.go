package toolbox

import (
	"testing"
)

var _cities []City

func init() {
	cs, err := DeserializationToCities("../tmp.json")
	if err != nil {
		panic(err)
	}

	_cities = cs
}

func TestCheckCoordinates(t *testing.T) {
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates(c, _cities)
	t.Log(parent, code)
}

func BenchmarkCheckCoordinates(b *testing.B) {
	b.ResetTimer()
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates(c, _cities)
	b.Log(parent, code)
}
