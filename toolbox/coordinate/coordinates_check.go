package coordinate

import (
	"sort"
)

// CheckCoordinates 坐标反查，先查找城市再查找省
func CheckCoordinates(c XY, cities []City, cityCenters []Center) (parent, code string) {
	maybeCities := Cities{}
	for i, v := range cities {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {
			cities[i].Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
			maybeCities = append(maybeCities, &cities[i])
		}
	}
	sort.Sort(maybeCities)
	find := false
	maybeProvinces := make(map[string]*Province)
	for i := range maybeCities {
		maybeProvinces[maybeCities[i].Parent] = maybeCities[i].Province
		for _, plg := range maybeCities[i].Polygons {
			if c.In(plg) {
				find = true
				return maybeCities[i].Parent, maybeCities[i].Code
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

	// 返回距离城市中心最近的城市
	minDistance := 0.0
	city := City{}
	for _, v := range cityCenters {
		distance := geoDistance(c.X, c.Y, v.Coordinate.Lng, v.Coordinate.Lat)
		if minDistance == 0.0 || minDistance > distance {
			minDistance = distance
			city = *v.City
		}
	}
	if minDistance > 0.04936914266221848 {
		return "", ""
	}
	
	return city.Parent, city.Code
}

// CheckCoordinates2 坐标反查, 先查找省再查找城市
func CheckCoordinates2(c XY, provinces []Province, cityCenters []Center) (parent, code string) {
	maybeProvinces := Provinces{}
	for i, v := range provinces {
		if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
			c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

			provinces[i].Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
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

	if hitProvince != nil {
		maybeCities := Cities{}
		for ci, v := range hitProvince.Cites {
			if c.X <= v.Apices.MaxLng && c.X >= v.Apices.MinLng &&
				c.Y <= v.Apices.MaxLat && c.Y >= v.Apices.MinLat {

				hitProvince.Cites[ci].Distance = geoDistance(c.X, c.Y, v.Center.Lng, v.Center.Lat)
				maybeCities = append(maybeCities, &hitProvince.Cites[ci])
			}
		}

		sort.Sort(maybeCities)
		for _, city := range maybeCities {
			for _, plg := range city.Polygons {
				if c.In(plg) {
					code = city.Code
					parent = city.Parent
					return
				}
			}
		}
	}

	// 返回距离城市中心最近的城市
	minDistance := 0.0
	var city *City
	for _, v := range cityCenters {
		distance := geoDistance(c.X, c.Y, v.Coordinate.Lng, v.Coordinate.Lat)
		if minDistance == 0.0 || minDistance > distance {
			minDistance = distance
			city = v.City
		}
	}
	if minDistance > 0.04936914266221848 {
		return "", ""
	}

	return city.Parent, city.Code
}
