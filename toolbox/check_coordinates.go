package toolbox

import (
	"encoding/json"
	"math"
	"os"
	"sort"
)

// China ...
type China struct {
	Province []Province `json:"provice"`
}

// Province ...
type Province struct {
	Parent  string `json:"parent_code"`
	Code    string `json:"code"`
	Polygon []XY   `json:"-"`
	Cites   []City `json:"subordinate"`
}

// City ...
type City struct {
	Parent string `json:"parent_code"`
	Code   string `json:"code"`

	Coordinates [][2]float64 `json:"coordinates"`
	Polygon     []XY         `json:"-"`

	Apices Apices `json:"apices"`

	Center Coordinate `json:"center"`

	Distance float64 `json:"-"`
}

// Apices ...
type Apices struct {
	MaxLng float64 `json:"maxlng"`
	MinLng float64 `json:"minlng"`
	MaxLat float64 `json:"maxlat"`
	MinLat float64 `json:"minlat"`
}

// Coordinate ...
type Coordinate struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// DeserializationToCities ...
func DeserializationToCities(path string) ([]City, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	china := China{}
	err = json.NewDecoder(f).Decode(&china)
	if err != nil {
		return nil, err
	}

	cities := []City{}
	for _, ps := range china.Province {
		for _, city := range ps.Cites {
			for _, coor := range city.Coordinates {
				city.Polygon = append(city.Polygon, XY{coor[0], coor[1]})
			}
			cities = append(cities, city)
		}
	}

	return cities, nil
}

// Cities ...
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

// geoDistance 返回不是准确的距离，是一个相对值
func geoDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
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
	// 不需要计算精确值
	// dist = dist * 180 / math.Pi
	// dist = dist * 60 * 1.1515

	return dist
}

// CheckCoordinates 坐标反查
func CheckCoordinates(c XY, cities []City) (parent, code string) {
	maybeCities := Cities{}
	for i, v := range cities {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			city := cities[i]
			maybeCities = append(maybeCities, &city)
		}
	}

	sort.Sort(maybeCities)
	for _, v := range maybeCities {
		if c.In(v.Polygon) {
			return v.Parent, v.Code
		}
	}

	return "", ""
}
