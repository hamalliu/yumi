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

	Attributes map[string]string `json:"attributes"`
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
		attrs.WithValues(AttrKeyColor, ins.Attributes[AttrKeyColor])
		attrs.WithValues(AttrKeyWeight, ins.Attributes[AttrKeyWeight])
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
