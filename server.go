package imosrpc

import (
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var httpPort = flag.String(
	"imosrpc-http", "", "Address to bind for ImosRPC server.")
var httpListener net.Listener

func Serve() {
	var wg sync.WaitGroup

	flag.Parse()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if *httpPort != "" || httpListener != nil {
			var err error
			if httpListener == nil {
				httpListener, err = net.Listen("tcp", *httpPort)
			}
			if err != nil {
				panic(err)
			}
			log.Printf("listening http://%s/...\n", httpListener.Addr())
			server := http.Server{
				Addr:           httpListener.Addr().String(),
				Handler:        http.HandlerFunc(DefaultHandler),
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
			if err := server.Serve(httpListener); err != nil {
				panic(err)
			}
		}
	}()

	wg.Wait()
}
