package header

import (
	"net/http"
	"strings"
)

const (
	// http head
	_httpHeaderRequestURI      = "x-yumi-request-uri"
	_httpHeaderContentSecurity = "x-yumi-content-security"
	_httpHeaderRemoteIP        = "x-yumi-real-ip"
	_httpHeaderRemotePort      = "x-yumi-real-port"
)

//ReqHeaders ...
func ReqHeaders() []string {
	return []string{_httpHeaderRemoteIP, _httpHeaderRemotePort}
}

//RespHeaders ...
func RespHeaders() []string {
	return nil
}

// RemoteIP implements a best effort algorithm to return the real client IP, it parses
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

//RemotePort ...
func RemotePort(req *http.Request) (port string) {
	if port = req.Header.Get(_httpHeaderRemotePort); port != "" && port != "null" {
		return
	}
	return
}

// ContentSecurity ...
func ContentSecurity(req *http.Request) (cs string) {
	if cs = req.Header.Get(_httpHeaderContentSecurity); cs != "" && cs != "null" {
		return
	}
	return
}

// RequestURI ...
func RequestURI(req *http.Request) (uri string) {
	if uri = req.Header.Get(_httpHeaderRequestURI); uri != "" && uri != "null" {
		return
	}
	return
}

