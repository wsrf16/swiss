package main

import (
	"github.com/wsrf16/swiss/sugar/logo"
)

func main() {
	// logo.SetFormatter(&logo.JSONFormatter{PrettyPrint: true})
	logo.SetFormatter(&logo.TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})
	logo.SetLevel(logo.DebugLevel)

	// {
	//    transfer, err := httptrans.BuildTransfer(":8082", false)
	//    if err != nil {
	//        panic(err)
	//    }
	//
	//    go func() {
	//        time.Sleep(10 * time.Second)
	//        transfer.Unload()
	//    }()
	//
	//    err = transfer.TransferFromListen()
	//    if err != nil {
	//        panic(err)
	//    }
	//    os.Exit(0)
	// }

	//cname, _ := net.LookupCNAME("www.baidu.com")
	//fmt.Println("cname:", cname)
	//dnsname, _ := net.LookupAddr("127.0.0.1")
	//fmt.Println("hostname:", dnsname)
	// fmt.Println(err)
	// select {}

}
