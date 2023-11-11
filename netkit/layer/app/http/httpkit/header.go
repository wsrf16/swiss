package httpkit

import (
	"net"
	"net/http"
)

func GetValue(r *http.Request, key string) ([]string, bool) {
	prior, ok := r.Header[key]
	return prior, ok
}

func AddHeader(header http.Header, key string, values []string) http.Header {
	for _, v := range values {
		header.Add(key, v)
	}
	return header
}

func GetXRealIP(r *http.Request) ([]string, bool) {
	return GetValue(r, "X-Real-IP")
}

func GetRemoteIP(r *http.Request) ([]string, error) {
	if ip, b := GetValue(r, "X-Real-IP"); b {
		return ip, nil
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		return nil, err
	} else {
		ip := make([]string, 1)
		ip[0] = host
		return ip, nil
	}
}

func GetXForwardedFor(r *http.Request) ([]string, bool) {
	return GetValue(r, "X-Forwarded-For")
}

func GetAddOnXForwardedFor(r *http.Request) ([]string, bool) {
	if prior, ok := GetXForwardedFor(r); ok {
		if realIP, ok := GetXRealIP(r); ok {
			n := append(prior, realIP...)
			return n, ok
		} else {
			n := append(prior, "unknown")
			return n, ok
		}
	} else {
		return GetXRealIP(r)
	}
}

func AddXForwardedFor(r *http.Request) {
	forwardedFor, _ := GetAddOnXForwardedFor(r)
	r.Header["X-Forwarded-For"] = forwardedFor
}

//func AddXForwarderFoxAndXRealIP(tarReq *http.Request) {
//    AddXRealIP(tarReq)
//    AddXForwardedFor(tarReq)
//}
