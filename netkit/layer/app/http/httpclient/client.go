package httpclient

import (
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"os"
	"strings"
)

func marshal(data interface{}) (string, error) {
	return jsonkit.MarshalToJson(data)
}

func Request(method, url string, heads map[string]string, body io.Reader) (*http.Response, error) {
	if request, err := http.NewRequest(method, url, body); err != nil {
		return nil, err
	} else {
		if heads != nil {
			for k, v := range heads {
				request.Header.Add(k, v)
			}
		}
		return http.DefaultClient.Do(request)
	}
}

func RequestBasicAuth(method, url string, body io.Reader, username string, password string) (*http.Response, error) {
	if request, err := http.NewRequest(method, url, body); err != nil {
		return nil, err
	} else {
		if len(username) == 0 && len(password) == 0 {
			request.SetBasicAuth(username, password)
		}
		return http.DefaultClient.Do(request)
	}
}

func RequestT[T any](method, url string, t T, heads map[string]string) (*http.Response, error) {
	bodyMal, err := marshal(t)
	if err != nil {
		return nil, err
	}
	return Request(method, url, heads, strings.NewReader(bodyMal))
}

// func PostFormT[T any](url string, t T, heads map[string]string) (*http.Response, error) {
//     bodyMal, err := marshal(t)
//     if err != nil {
//         return nil, err
//     }
//
//     return http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(bodyMal))
// }
//
// func PostForm(url string, data url.Values) (*http.Response, error) {
//     return http.PostForm(url, data)
// }

func Socks5ProxyClient(network, address string, auth *proxy.Auth, forward proxy.Dialer) (*http.Client, error) {
	dialer, err := proxy.SOCKS5(network, address, auth, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient, nil
}

func httpProxyClient(address string) *http.Client {
	os.Setenv("HTTP_PROXY", address)
	os.Setenv("HTTPS_PROXY", address)
	c := &http.Client{}
	return c
}
