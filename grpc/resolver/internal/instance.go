package internal

import (
	"net/url"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

const (
	AttrKeyRegion  = "region"
	AttrKeyZone    = "zone"
	AttrKeyCluster = "cluster"
	AttrKeyColor   = "color"
	AttrKeyWeight  = "weight"
)

// Instance represents a server the client connects to.
type Instance struct {
	// Region is region.
	Region string `json:"region"`
	// Zone is IDC.
	Zone string `json:"zone"`
	// Cluster ...
	Cluster string `json:"cluster"`
	// Color ...
	Color string `json:"color"`
	// Weight ...
	Weight int `json:"weight"`

	// AppID is mapping servicetree appid.
	AppID string `json:"appid"`
	// Env prod/pre、uat/fat1
	Env string `json:"env"`
	// Hostname is hostname from docker.
	Hostname string `json:"hostname"`
	// Addrs is the address of app instance
	// format: scheme://host
	Addrs []string `json:"addrs"`
	// Version is publishing version.
	Version string `json:"version"`
	// LastTs is instance latest updated timestamp
	LastTs int64 `json:"latest_timestamp"`
	// Status instance status, eg: 1UP 2Waiting
	Status int64 `json:"status"`
}

// ToGrpcAddress ...
func ToGrpcAddress(inss []*Instance) []resolver.Address {
	addrs := []resolver.Address{}
	for _, ins := range inss {
		addr := resolver.Address{}
		addr.Type = resolver.Backend
		// addr.ServerName = ins.AppID
		for _, a := range ins.Addrs {
			u, err := url.Parse(a)
			if err == nil && u.Scheme == "grpc" {
				addr.Addr = u.Host
			}
		}
		attrs := attributes.Attributes{}
		attrs.WithValues(AttrKeyColor, ins.Color)
		attrs.WithValues(AttrKeyWeight, ins.Weight)
		addr.Attributes = &attrs

		addrs = append(addrs, addr)
	}

	return addrs
}

// InstancesInfo instance info.
type InstancesInfo struct {
	Instances map[string][]*Instance `json:"instances"`
	LastTs    int64                  `json:"latest_timestamp"`
	Scheduler *Scheduler             `json:"scheduler"`
}

// Scheduler scheduler.
type Scheduler struct {
	Clients map[string]*ZoneStrategy `json:"clients"`
}

// ZoneStrategy is the scheduling strategy of all zones
type ZoneStrategy struct {
	Zones map[string]*Strategy `json:"zones"`
}

// Strategy is zone scheduling strategy.
type Strategy struct {
	Weight int64 `json:"weight"`
}
