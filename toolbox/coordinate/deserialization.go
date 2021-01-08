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

	Coordinates [][][2]float64 `json:"coordinates"`
	Polygons    [][]XY         `json:"-"`

	Apices Apices `json:"apices"`

	Center Coordinate `json:"center"`

	Cites []City `json:"subordinate"`

	Distance float64 `json:"-"`
}

// City ...
type City struct {
	Parent string `json:"parent_code"`
	Code   string `json:"code"`

	Coordinates [][][2]float64 `json:"coordinates"`
	Polygons    [][]XY         `json:"-"`

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

// Center ...
type Center struct {
	Coordinate Coordinate
	City *City
}

// DeserializationToCities ...
func DeserializationToCities(path string) ([]City, []Center, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	china := China{}
	err = json.NewDecoder(f).Decode(&china)
	if err != nil {
		return nil, nil, err
	}

	cities := []City{}
	cityCenters := []Center{}
	for pi := range china.Province {
		for _, plg := range china.Province[pi].Coordinates {
			polygon := Polygon{}
			for _, coor := range plg {
				polygon = append(polygon, XY{coor[0], coor[1]})
			}
			china.Province[pi].Polygons = append(china.Province[pi].Polygons, polygon)
		}

		for ci := range china.Province[pi].Cites {
			for _, plg := range china.Province[pi].Cites[ci].Coordinates {
				polygon := Polygon{}
				for _, coor := range plg {
					polygon = append(polygon, XY{coor[0], coor[1]})
				}
				china.Province[pi].Cites[ci].Polygons = append(china.Province[pi].Cites[ci].Polygons, polygon)
			}
			china.Province[pi].Cites[ci].Province = &china.Province[pi]
			cities = append(cities, china.Province[pi].Cites[ci])

			cityCenters = append(cityCenters, Center{Coordinate: china.Province[pi].Cites[ci].Center, City: &china.Province[pi].Cites[ci]})
		}
	}

	return cities, cityCenters, nil
}

// DeserializationToProvince ...
func DeserializationToProvince(path string) ([]Province, []Center, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	china := China{}
	err = json.NewDecoder(f).Decode(&china)
	if err != nil {
		return nil, nil, err
	}

	provinces := []Province{}
	cityCenters := []Center{}
	for _, ps := range china.Province {
		for _, plg := range ps.Coordinates {
			polygon := Polygon{}
			for _, coor := range plg {
				polygon = append(polygon, XY{coor[0], coor[1]})
			}
			ps.Polygons = append(ps.Polygons, polygon)
		}

		for ci := range ps.Cites {
			for _, plg := range ps.Cites[ci].Coordinates {
				polygon := Polygon{}
				for _, coor := range plg {
					polygon = append(polygon, XY{coor[0], coor[1]})
				}
				ps.Cites[ci].Polygons = append(ps.Cites[ci].Polygons, polygon)
			}

			cityCenters = append(cityCenters, Center{Coordinate: ps.Cites[ci].Center, City: &ps.Cites[ci]})
		}

		provinces = append(provinces, ps)
	}

	return provinces, cityCenters, nil
}
