package coordinate

import (
	"testing"
)

var _cities []City
var _provinces []Province
var _cityCenters []Center

func init() {
	cs, cityCenters, err := DeserializationToCities("./行政区划.json")
	if err != nil {
		panic(err)
	}

	_cities = cs
	_cityCenters = cityCenters

	ps, _, err := DeserializationToProvince("./行政区划.json")
	if err != nil {
		panic(err)
	}

	_provinces = ps
}

func TestCheckCoordinates(t *testing.T) {
	c := XY{121.778859, 31.310198}
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
	c := XY{121.778859, 31.310198}
	parent, code := CheckCoordinates2(c, _provinces)
	t.Log(parent, code)
}

func BenchmarkCheckCoordinates2(b *testing.B) {
	b.ResetTimer()
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates2(c, _provinces)
	b.Log(parent, code)
}
