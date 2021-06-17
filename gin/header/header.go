package header

import (
	"net/http"
	"strings"
)

const (
	// http head
	_httpHeaderRemoteIP        = "x-yumi-real-ip"
	_httpHeaderRemotePort      = "x-yumi-real-port"
	_httpHeaderContentSecurity = "x-yumi-content-security"
	_httpHeaderBodyEncrypt     = "x-yumi-body-encrypt"
)

//ReqHeaders ...
func ReqHeaders() []string {
	return []string{_httpHeaderRemoteIP, _httpHeaderRemotePort, _httpHeaderBodyEncrypt, _httpHeaderContentSecurity}
}

//RespHeaders ...
func RespHeaders() []string {
	return []string{_httpHeaderBodyEncrypt}
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

//RemotePort 获取客户端端口
func RemotePort(req *http.Request) (port string) {
	if port = req.Header.Get(_httpHeaderRemotePort); port != "" && port != "null" {
		return
	}
	return
}

// ContentSecurity 获取安全传输校验数据
func ContentSecurity(req *http.Request) (security string) {
	return req.Header.Get(_httpHeaderContentSecurity)
}

// BodyEncrypt 获取body是否加密
func BodyEncrypt(req *http.Request) (security string) {
	return req.Header.Get(_httpHeaderBodyEncrypt)
}

// SetBodyEncrypt 设置body是否加密为true
func SetBodyEncrypt(w http.ResponseWriter) {
	w.Header().Set(_httpHeaderBodyEncrypt, "true")
}
