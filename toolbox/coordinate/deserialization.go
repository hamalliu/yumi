package coordinate

import (
	"encoding/json"
	"os"
)

// China ...
type China struct {
	Province []Province `json:"provice"`
}

// Province ...
type Province struct {
	Parent string `json:"parent_code"`
	Code   string `json:"code"`

	Coordinates [][2]float64 `json:"coordinates"`
	Polygon     []XY         `json:"-"`

	Apices Apices `json:"apices"`

	Center Coordinate `json:"center"`

	Cites []City `json:"subordinate"`

	Distance float64 `json:"-"`
}

// City ...
type City struct {
	Parent string `json:"parent_code"`
	Code   string `json:"code"`

	Coordinates [][2]float64 `json:"coordinates"`
	Polygon     []XY         `json:"-"`

	Apices Apices `json:"apices"`

	Center Coordinate `json:"center"`

	Province *Province `json:"-"`
	Distance float64   `json:"-"`
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
	for pi, ps := range china.Province {
		for _, coor := range china.Province[pi].Coordinates {
			china.Province[pi].Polygon = append(china.Province[pi].Polygon, XY{coor[0], coor[1]})
		}

		for _, city := range ps.Cites {
			for _, coor := range city.Coordinates {
				city.Polygon = append(city.Polygon, XY{coor[0], coor[1]})
			}
			city.Province = &china.Province[pi]
			cities = append(cities, city)
		}
	}

	return cities, nil
}

// DeserializationToProvince ...
func DeserializationToProvince(path string) ([]Province, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	china := China{}
	err = json.NewDecoder(f).Decode(&china)
	if err != nil {
		return nil, err
	}

	provinces := []Province{}
	for _, ps := range china.Province {
		for _, coor := range ps.Coordinates {
			ps.Polygon = append(ps.Polygon, XY{coor[0], coor[1]})
		}

		for ci := range ps.Cites {
			for _, coor := range ps.Cites[ci].Coordinates {
				ps.Cites[ci].Polygon = append(ps.Cites[ci].Polygon, XY{coor[0], coor[1]})
			}
		}

		provinces = append(provinces, ps)
	}

	return provinces, nil
}
