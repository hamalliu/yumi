package toolbox

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

type China struct {
	Provinces []Province
}

type Province struct {
	Parent    string
	Code      string
	Rectangle []XY
	Polygon   []XY
	Cites     []City
}

type City struct {
	Parent    string
	Code      string
	Rectangle []XY
	Center    XY
	Polygon   []XY

	Distance float64
}

type Cities []*City

func (c Cities) Len() int {
    return len(c)
}
func (c Cities) Less(i, j int) bool {
    return c[i].Distance < c[j].Distance
}
func (c Cities) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
}

func GeoDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
	// const PI float64 = 3.141592653589793

	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	return dist
}

func TestRayCast(t *testing.T) {
	cities := []City{}
	coordinate := XY{116.16472347941438, 31.127336498964578}
	maybeCities := Cities{}
	for _, v := range cities {
		if coordinate.In(v.Rectangle) {
			v.Distance = GeoDistance(coordinate.X, coordinate.Y, v.Center.X, v.Center.Y)
			maybeCities = append(maybeCities, &v)
		}
	}

	sort.Sort(maybeCities)
	for _, v := range maybeCities {
		if coordinate.In(v.Polygon) {
			fmt.Println(v.Code, v.Parent)
		}
	}
	t.Log("未定位成功")
}
