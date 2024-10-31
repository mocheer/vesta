package vesta

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/chromedp/cdproto/fetch"
)

// FetchByHttp
// 这是通过http重新请求
func FetchByHttp(e *fetch.EventRequestPaused) []byte {
	request := e.Request
	client := &http.Client{}
	// 创建请求
	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		panic(err)
	}
	// 添加请求头
	for key, val := range request.Headers {
		req.Header.Add(key, val.(string))
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 使用 resp.Body 读取响应内容
	body, _ := io.ReadAll(resp.Body)
	return body
}

// UpdateFetchBody
func UpdateFetchBody(ctx context.Context, e *fetch.EventRequestPaused, callback func([]byte, *http.Response) []byte) {
	request := e.Request
	client := &http.Client{}
	// 创建请求
	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		panic(err)
	}
	// 添加请求头
	for key, val := range request.Headers {
		req.Header.Add(key, val.(string))
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 使用 resp.Body 读取响应内容
	body, _ := io.ReadAll(resp.Body)
	body = callback(body, resp)

	// 可能需要根据最新的body更新Content-Length
	headers := make([]*fetch.HeaderEntry, 0)
	for key, val := range resp.Header {
		headers = append(headers, &fetch.HeaderEntry{Name: key, Value: val[0]})
	}
	//

	fetch.FulfillRequest(e.RequestID, int64(resp.StatusCode)).
		WithBody(base64.StdEncoding.EncodeToString([]byte(body))).
		WithResponseHeaders(headers).
		WithResponsePhrase("OK").
		Do(ctx)
}

// ReplaceFetchBody
func ReplaceFetchBody(ctx context.Context, e *fetch.EventRequestPaused, body string) {
	headers := make([]*fetch.HeaderEntry, 0)
	headers = append(headers, &fetch.HeaderEntry{Name: "Connection", Value: "closed"})
	headers = append(headers, &fetch.HeaderEntry{Name: "Content-Length", Value: "6"})
	headers = append(headers, &fetch.HeaderEntry{Name: "Content-Type", Value: "text/javascript"})
	fetch.FulfillRequest(e.RequestID, 200).
		WithBody(base64.StdEncoding.EncodeToString([]byte(body))).
		WithResponseHeaders(headers).
		WithResponsePhrase("OK").
		Do(ctx)
}
