package coordinate

import (
	"sort"
)

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
	find := false
	maybeProvinces := make(map[string]*Province)
	for _, city := range maybeCities {
		maybeProvinces[city.Parent] = city.Province
		if c.In(city.Polygon) {
			find = true
			return city.Parent, city.Code
		}
	}
	if !find {
		for _, province := range maybeProvinces {
			if c.In(province.Polygon) {
				find = true
				return province.Code, ""
			}
		}
	}

	return "", ""
}

// CheckCoordinates2 坐标反查
func CheckCoordinates2(c XY, provinces []Province) (parent, code string) {
	maybeProvinces := Provinces{}
	for i, v := range provinces {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			province := provinces[i]
			maybeProvinces = append(maybeProvinces, &province)
		}
	}

	sort.Sort(maybeProvinces)
	var hitProvince *Province
	for _, province := range maybeProvinces {
		if c.In(province.Polygon) {
			parent = province.Code
			hitProvince = province
		}
	}

	maybeCities := Cities{}
	for ci, v := range hitProvince.Cites {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			city := hitProvince.Cites[ci]
			maybeCities = append(maybeCities, &city)
		}
	}

	sort.Sort(maybeCities)
	for _, city := range maybeCities {
		if c.In(city.Polygon) {
			code = city.Code
		}
	}

	return
}
