package main

import (
	"fmt"
	"github.com/wsrf16/swiss/module/mq/kafka"
	"github.com/wsrf16/swiss/sugar/logo"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Item struct {
	A int
}

func main() {
	items := make([]Item, 2, 10)
	items[0] = Item{}
	items[0].A = 0
	items[1] = Item{}
	items[1].A = 1

	var all []*Item
	for _, item := range items {
		item := item
		all = append(all, &item)
	}
	fmt.Println(all)

	// logo.SetFormatter(&logo.JSONFormatter{PrettyPrint: true})
	// // cname, _ := net.LookupCNAME("www.baidu.com")
	// // fmt.Println("cname:", cname)
	// // dnsname, _ := net.LookupAddr("127.0.0.1")
	// // fmt.Println("hostname:", dnsname)
	// err := sockskit.TransferFromListenAddress(":1080")
	// fmt.Println(err)
	// select {}

	producer, err := kafka.NewSyncProducer([]string{"mecs.com:19092"}, "", "")
	if err != nil {
		logo.Error("a", err)
	}

	obj := make(map[string]interface{})
	obj["name"] = time.Now().String()
	obj["age"] = rand.Intn(100)

	err = kafka.SendT(obj, "application-log", producer)
	if err != nil {
		logo.Error("b", err)
	}

}

func handleHttps(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
