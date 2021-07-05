package coordinate

import "math"

// Provinces ...
type Provinces []*Province

func (p Provinces) Len() int {
	return len(p)
}
// 从小到大
func (p Provinces) Less(i, j int) bool {
	return p[i].Distance < p[j].Distance
}
func (p Provinces) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Cities ...
type Cities []*City

func (c Cities) Len() int {
	return len(c)
}
// 从小到大
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
