package coordinate

import (
	"fmt"
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
	ps := []XY{
		{121.385531, 30.985695},
		{121.381586, 30.984675},
		{121.375223, 30.982826},
		{121.359979, 30.97816},
		{121.778859, 31.310198},
		{121.35871, 30.977858},
	}
	for _, v := range ps {
		t.Log(v)
		parent, code := CheckCoordinates(v, _cities, _cityCenters)
		t.Log(parent, code)
	}
}

func BenchmarkCheckCoordinates(b *testing.B) {
	b.ResetTimer()
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates(c, _cities, _cityCenters)
	b.Log(parent, code)
}

func TestCheckCoordinates2(t *testing.T) {
	ps := []XY{
		{121.385531, 30.985695},
		{121.381586, 30.984675},
		{121.375223, 30.982826},
		{121.359979, 30.97816},
		{121.778859, 31.310198},
		{121.35871, 30.977858},
		{121.412513, 31.1892},
		{61.180859,51.222905},
	}
	for _, v := range ps {
		t.Log(v)
		parent, code := CheckCoordinates2(v, _provinces, _cityCenters)
		t.Log(parent, code)
	}
}

func BenchmarkCheckCoordinates2(b *testing.B) {
	b.ResetTimer()
	c := XY{116.16472347941438, 31.127336498964578}
	parent, code := CheckCoordinates2(c, _provinces, _cityCenters)
	b.Log(parent, code)
}

func TestDistance(t *testing.T) {
	fmt.Println(geoDistance(121.472644, 31.231706, 124.132519, 29.576592))
}
