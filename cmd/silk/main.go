package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elazarl/goproxy"
)

type Meta struct {
	req      *http.Request
	resp     *http.Response
	err      error
	t        time.Time
	sess     int64
	bodyPath string
	from     string
}

var emptyResp = &http.Response{}
var emptyReq = &http.Request{}

func StartProxy() {
	log.Println("Starting proxy server")
	go func() {
		proxy := goproxy.NewProxyHttpServer()

		// var handleFunc goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		// 	return goproxy.OkConnect, host
		// }
		// proxy.OnRequest().HandleConnect(handleFunc)
		// proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
		// proxy.OnRequest().HandleConnect(goproxy.AlwaysReject)
		proxy.AllowHTTP2 = true
		proxy.Verbose = true
		//
		// proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		// 	from := ""
		// 	if ctx.UserData != nil {
		// 		from = ctx.UserData.(*transport.RoundTripDetails).TCPAddr.String()
		// 	}
		// 	if resp == nil {
		// 		resp = emptyResp
		// 	} else {
		// 		var buf bytes.Buffer
		// 		tee := io.TeeReader(resp.Body, &buf)
		// 		resp.Body = io.NopCloser(tee)
		// 		log.Println("-----", len(buf.Bytes()))
		// 	}
		// 	m := &Meta{
		// 		resp: resp,
		// 		err:  ctx.Error,
		// 		t:    time.Now(),
		// 		sess: ctx.Session,
		// 		from: from}
		// 	//
		// 	log.Println("-----", *m)
		// 	return resp
		// })
		if err := http.ListenAndServe(":9080", proxy); err != nil {
			log.Fatal(err)
		}
	}()
}

func main() {
	StartProxy()
	select {}
}
