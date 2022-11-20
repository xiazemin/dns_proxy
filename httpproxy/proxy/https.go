package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func ServeHttps(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 60*time.Second)
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

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	fmt.Println(r.RequestURI, "--->", destConn.RemoteAddr().String())

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

	buf, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(buf))
	bufr, _ := httputil.DumpResponse(&http.Response{}, true)
	fmt.Println(string(bufr))
	// dump, wait := destConn.Dump()
	//https://github.com/buger/goreplay/blob/master/input_http.go
	// https://github.com/buger/goreplay/blob/f3bb50192cf2c5a95627ac40e151ffd5b6fdd719/output_http.go

	bytesData := make([]byte, 1024)
	middle := bytes.NewBuffer(bytesData)
	io.Copy(middle, destConn)
	fmt.Println(middle.String())
	middle.Reset()
	fmt.Println(clientConn.Read(bytesData))
	fmt.Println(string(bytesData))
	//会卡死

}

//https://vimsky.com/examples/detail/golang-ex-github.com.rhino1998.god.client-Conn-Dump-method.html
// dump, wait := conn.Dump()
