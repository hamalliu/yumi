package coordinate

import (
	"testing"
)

var _cities []City
var _provinces []Province

func init() {
	cs, err := DeserializationToCities("./division.json")
	if err != nil {
		panic(err)
	}

	_cities = cs

	ps, err := DeserializationToProvince("./division.json")
	if err != nil {
		panic(err)
	}

	_provinces = ps
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

func TestCheckCoordinates2(t *testing.T) {
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates2(c, _provinces)
	t.Log(parent, code)
}

func BenchmarkCheckCoordinates2(b *testing.B) {
	b.ResetTimer()
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates2(c, _provinces)
	b.Log(parent, code)
}
