package types

import (
	"errors"
	"net"
)

type IP net.IP

// UnmarshalJSON ...
func (ip *IP) UnmarshalJSON(data []byte) (err error) {
	netIP := net.ParseIP(string(data[1 : len(data)-1]))
	if netIP == nil {
		return errors.New("data: " + string(data) + "not in ip format")
	}
	*ip = IP(netIP.To16())
	return nil
}

// MarshalJSON ...
func (ip *IP) MarshalJSON() ([]byte, error) {
	intBytes := []byte(*ip)
	var ret net.IP
	if len(intBytes) == net.IPv4len {
		ret = net.IPv4(intBytes[0], intBytes[1], intBytes[2], intBytes[3])
	}
	if len(intBytes) == net.IPv6len {
		ret = net.IP(intBytes)
	}
	
	retStr := ret.String()
	return []byte(`"` + retStr + `"`), nil
}
