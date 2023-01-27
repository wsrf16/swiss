package httpkit

import (
	"net"
	"net/http"
)

func GetXRealIP(r *http.Request) (string, error) {
	if ip := r.Header.Get("X-Real-IP"); len(ip) != 0 {
		return ip, nil
	}
	// http.DefaultTransport.RoundTrip()
	// http.DefaultClient
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		return "", err
	} else {
		return ip, nil
	}
}

func GetXForwardedFor(r *http.Request) ([]string, bool) {
	if prior, ok := r.Header["X-Forwarded-For"]; ok {
		return prior, ok
	} else {
		return nil, false
	}
}

func GetAddOnXForwardedFor(r *http.Request) ([]string, bool) {
	if prior, ok := GetXForwardedFor(r); ok {
		if realIP, err := GetXRealIP(r); err != nil {
			n := append(prior, "unknown")
			return n, false
		} else {
			n := append(prior, realIP)
			return n, true
		}
	} else {
		if realIP, err := GetXRealIP(r); err != nil {
			return nil, false
		} else {
			return []string{realIP}, true
		}
	}
}

func AddXRealIP(r *http.Request) error {
	realIP, err := GetXRealIP(r)
	if err != nil {
		return err
	} else {
		r.Header.Set("X-Real-IP", realIP)
		return nil
	}
}

func AddXForwardedFor(r *http.Request) {
	forwardedFor, _ := GetAddOnXForwardedFor(r)
	r.Header["X-Forwarded-For"] = forwardedFor
}

func AddXForwarderFoxAndXRealIP(tarReq *http.Request) {
	AddXRealIP(tarReq)
	AddXForwardedFor(tarReq)
}
