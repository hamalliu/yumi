package header

import (
	"net/http"
	"strings"
)

const (
	// http head
	_httpHeaderUserId     = "x-yumi-user-id"
	_httpHeaderRemoteIP   = "x-yumi-real-ip"
	_httpHeaderRemotePort = "x-yumi-real-port"
)

func ReqHeaders() []string {
	return []string{_httpHeaderUserId, _httpHeaderRemoteIP, _httpHeaderRemotePort}
}

func RespHeaders() []string {
	return nil
}

func UserId(req *http.Request) string {
	return req.Header.Get(_httpHeaderUserId)
}

// remoteIP implements a best effort algorithm to return the real client IP, it parses
// x-backend-bm-real-ip or X-Real-IP or X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func RemoteIP(req *http.Request) (remote string) {
	if remote = req.Header.Get(_httpHeaderRemoteIP); remote != "" && remote != "null" {
		return
	}
	var xff = req.Header.Get("X-Forwarded-For")
	if idx := strings.IndexByte(xff, ','); idx > -1 {
		if remote = strings.TrimSpace(xff[:idx]); remote != "" {
			return
		}
	}
	if remote = req.Header.Get("X-Real-IP"); remote != "" {
		return
	}
	remote = req.RemoteAddr[:strings.Index(req.RemoteAddr, ":")]
	return
}

func RemotePort(req *http.Request) (port string) {
	if port = req.Header.Get(_httpHeaderRemotePort); port != "" && port != "null" {
		return
	}
	return
}
