package coordinate

import (
	"sort"
)

// CheckCoordinates 坐标反查，先查找城市再查找省
func CheckCoordinates(c XY, cities []City) (parent, code string) {
	maybeCities := Cities{}
	for i, v := range cities {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			maybeCities = append(maybeCities, &cities[i])
		}
	}

	sort.Sort(maybeCities)
	find := false
	maybeProvinces := make(map[string]*Province)
	for _, city := range maybeCities {
		maybeProvinces[city.Parent] = city.Province
		for _, plg := range city.Polygons {
			if c.In(plg) {
				find = true
				return city.Parent, city.Code
			}
		}
	}
	if !find {
		for _, province := range maybeProvinces {
			for _, plg := range province.Polygons {
				if c.In(plg) {
					find = true
					return province.Code, ""
				}
			}
		}
	}

	return "", ""
}

// CheckCoordinates2 坐标反查, 先查找省再查找城市
func CheckCoordinates2(c XY, provinces []Province) (parent, code string) {
	maybeProvinces := Provinces{}
	for i, v := range provinces {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			maybeProvinces = append(maybeProvinces, &provinces[i])
		}
	}

	sort.Sort(maybeProvinces)
	var hitProvince *Province
	for _, province := range maybeProvinces {
	outprovince:
		for _, plg := range province.Polygons {
			if c.In(plg) {
				parent = province.Code
				hitProvince = province
				break outprovince
			}
		}
	}

	if hitProvince == nil {
		// 未定位成功
		return
	}

	maybeCities := Cities{}
	for ci, v := range hitProvince.Cites {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			v.Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			maybeCities = append(maybeCities, &hitProvince.Cites[ci])
		}
	}

	sort.Sort(maybeCities)
	for _, city := range maybeCities {
	outcity:
		for _, plg := range city.Polygons {
			if c.In(plg) {
				code = city.Code
				break outcity
			}
		}
	}

	return
}
