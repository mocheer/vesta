package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ouqiang/goproxy"
	"github.com/ouqiang/websocket"
)

// 实现证书缓存接口
type Cache struct {
	m sync.Map
}

func (c *Cache) Set(host string, cert *tls.Certificate) {
	c.m.Store(host, cert)
}
func (c *Cache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}

func main() {
	//goproxy.WithDecryptHTTPS(&Cache{}),
	proxy := goproxy.New(goproxy.WithDelegate(&EventHandler{}), goproxy.WithEnableWebsocketIntercept())
	server := &http.Server{
		Addr:         ":9080",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	log.Println("start listen 9080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type EventHandler struct{}

func (e *EventHandler) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {
	// 保存的数据可以在后面的回调方法中获取
	// ctx.Data["req_id"] = "uuid"

	// 禁止访问某个域名
	// if strings.Contains(ctx.Req.URL.Host, "example.com") {
	// 	rw.WriteHeader(http.StatusForbidden)
	// 	ctx.Abort()
	// 	return
	// }
}

func (e *EventHandler) Auth(ctx *goproxy.Context, rw http.ResponseWriter) {
	// log.Println("认证", ctx.Req.URL)
	// 身份验证
}

func (e *EventHandler) BeforeRequest(ctx *goproxy.Context) {
	log.Println("发送", ctx.Req.URL)
	if strings.Contains(ctx.Req.URL.String(), "online") {

		log.Println("----------------- 这个是一个websocket服务")
		log.Println(ctx.Req.URL.Scheme) // 这个可能是错的
		log.Println(ctx.Req.Method)
		log.Println(websocket.IsWebSocketUpgrade(ctx.Req))
	}
	// 修改header
	// ctx.Req.Header.Add("X-Request-Id", ctx.Data["req_id"].(string))
	// 设置X-Forwarded-For
	// if clientIP, _, err := net.SplitHostPort(ctx.Req.RemoteAddr); err == nil {
	// 	if prior, ok := ctx.Req.Header["X-Forwarded-For"]; ok {
	// 		clientIP = strings.Join(prior, ", ") + ", " + clientIP
	// 	}
	// 	ctx.Req.Header.Set("X-Forwarded-For", clientIP)
	// }
	// 读取Body
	// body, err := io.ReadAll(ctx.Req.Body)
	// if err != nil {
	// 	// 错误处理
	// 	return
	// }
	// Request.Body只能读取一次, 读取后必须再放回去
	// Response.Body同理
	// ctx.Req.Body = io.NopCloser(bytes.NewReader(body))
}

func (e *EventHandler) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err != nil {
		return
	}
	// 修改response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("返回值读取错误")
		// 错误处理
		return
	}
	log.Println("返回", ctx.Req.URL, string(body[:100]), len(body))
	resp.Body = io.NopCloser(bytes.NewReader(body))
}

// 设置上级代理
func (e *EventHandler) ParentProxy(req *http.Request) (*url.URL, error) {
	// return url.Parse("http://localhost:1087")
	return nil, nil
}

func (e *EventHandler) Finish(ctx *goproxy.Context) {
	log.Printf("请求结束 URL:%s\n", ctx.Req.URL)
}

// 记录错误日志
func (e *EventHandler) ErrorLog(err error) {
	log.Println("错误", err)
}

// WebSocketSendMessage websocket发送消息
func (e *EventHandler) WebSocketSendMessage(ctx *goproxy.Context, messageType *int, p *[]byte) {
	log.Println("-------------------- websocket send:", messageType, string(*p))
}

// WebSockerReceiveMessage websocket接收 消息
func (e *EventHandler) WebSocketReceiveMessage(ctx *goproxy.Context, messageType *int, p *[]byte) {
	log.Println("-------------------- websocket receive:", messageType, string(*p))
}
