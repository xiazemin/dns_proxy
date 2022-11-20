package httpproxy

import (
	"net/http"

	"github.com/xiazemin/dns_proxy/httpproxy/proxy"
)

func Serve() {

	server := &http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.Serve(w, r)
		}),
	}

	logger.Fatal(server.ListenAndServe())

}
