package discovery

import (
	"net/url"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"yumi/grpcwrapper/balancer"
)

// Instance represents a server the client connects to.
type Instance struct {
	// Region is region.
	Region string `json:"region"`
	// Zone is IDC.
	Zone string `json:"zone"`
	// Env prod/pre、uat/fat1
	Env string `json:"env"`

	// AppID is mapping servicetree appid.
	AppID string `json:"appid"`
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

	Balancer balancer.MetaFromResolver `json:"balancer"`

	Attributes map[string]interface{} `json:"attributes"`
}

// ToGrpcAddress ...
func ToGrpcAddress(inss []*Instance) []resolver.Address {
	addrs := []resolver.Address{}
	for _, ins := range inss {
		addr := resolver.Address{}
		addr.Type = resolver.Backend
		// 被（多证书）TLS,SSL服务器用于识别客户端需要的证书
		// addr.ServerName = ins.AppID
		for _, a := range ins.Addrs {
			u, err := url.Parse(a)
			if err == nil && u.Scheme == "grpc" {
				addr.Addr = u.Host
				continue
			}
		}
		attrs := attributes.Attributes{}
		attrs.WithValue(balancer.AttributesKey, ins.Balancer)
		addr.Attributes = &attrs

		addrs = append(addrs, addr)
	}

	return addrs
}
