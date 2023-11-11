package httpkit

import (
	"bytes"
	"crypto/tls"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func unmarshalSToT[T interface{}](js string) (*T, error) {
	return jsonkit.UnmarshalSToT[T](js)
}

func unmarshalBToT[T interface{}](bytes []byte) (*T, error) {
	return jsonkit.UnmarshalBToT[T](bytes)
}

func parseBody[T interface{}](r io.Reader) (*T, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return unmarshalBToT[T](bytes)
}

func ParseRequestBody[T interface{}](req *http.Request) (*T, error) {
	return parseBody[T](req.Body)
}

func ParseResponseBody[T any](resp *http.Response) (*T, error) {
	return parseBody[T](resp.Body)
}

func CloneRequest(lr *http.Request, url string) (*http.Request, error) {
	body, err := io.ReadAll(lr.Body)
	if err != nil {
		return nil, err
	}
	// r.Body = io.NopCloser(bytes.NewReader(body))

	rr, err := http.NewRequest(lr.Method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	rr.RemoteAddr = lr.RemoteAddr
	CopyHeaders(lr, rr)
	// rr.Header = lr.Header
	// rr.Host = lr.Host

	return rr, nil
}

func CopyHeaders(lr *http.Request, rr *http.Request) {
	rr.Header = make(http.Header)
	for k, v := range lr.Header {
		rr.Header[k] = v
	}
}

func WriteResponse(w http.ResponseWriter, resp *http.Response) {
	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func ForwardTo(lw http.ResponseWriter, lr *http.Request, url string) error {
	return ForwardToIntercept(lw, lr, url, nil, nil)
}

// func ForwardToHost(lw http.ResponseWriter, lr *http.Request) error {
//    url := urlkit.Join("http://"+lr.Host, lr.URL.String())
//    return ForwardToIntercept(lw, lr, url, nil, nil)
// }

func ForwardToIntercept(lw http.ResponseWriter, lr *http.Request, url string, rInterceptors []func(*http.Request), wInterceptors []func(http.ResponseWriter, *http.Response)) error {
	rr, err := CloneRequest(lr, url)
	defer rr.Body.Close()
	// if err != nil {
	//    return err
	// } else {
	//    AddXForwarderFoxAndXRealIP(rr)
	// }

	if rInterceptors != nil {
		for _, interceptor := range rInterceptors {
			interceptor(rr)
		}
	}

	res, err := new(http.Client).Do(rr)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	if wInterceptors != nil {
		for _, interceptor := range wInterceptors {
			interceptor(lw, res)
		}
	}
	// lw.Header().Set("Access-Control-Allow-Origin", "*")
	WriteResponse(lw, res)

	return nil
}

func ReverseProxy(lw http.ResponseWriter, lr *http.Request, uri string) error {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	})
	reverseProxy.ServeHTTP(lw, lr)
	return nil
}

func ParseSchema(req *http.Request) string {
	schema := lambda.If[string](req.TLS == nil, "http", "https")
	return schema
}

func ParseSchemaByTLS(tls *tls.ConnectionState) string {
	schema := lambda.If[string](tls == nil, "http", "https")
	return schema
}

func HeaderToMap(header http.Header) map[string]string {
	h := make(map[string]string)
	for k, _ := range header {
		v := header.Get(k)
		h[k] = v
	}
	return h
}
