package main

import (
	"github.com/wsrf16/swiss/sugar/logo"
)

func main() {
	logo.SetFormatter(&logo.JSONFormatter{PrettyPrint: true})
}
