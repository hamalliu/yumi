package internal

import (
	"math/rand"
	"sort"

	"github.com/dgryski/go-farm"
)

func subset(inss []*Instance, size int) []*Instance {
	backends := inss
	if len(backends) <= int(size) {
		return backends
	}
	clientID := ""
	sort.Slice(backends, func(i, j int) bool {
		return backends[i].Hostname < backends[j].Hostname
	})
	count := len(backends) / size
	// hash得到ID
	id := farm.Fingerprint64([]byte(clientID))
	// 获得rand轮数
	round := int64(id / uint64(count))

	s := rand.NewSource(round)
	ra := rand.New(s)
	//  根据source洗牌
	ra.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})
	start := (id % uint64(count)) * uint64(size)
	return backends[int(start) : int(start)+int(size)]
}
