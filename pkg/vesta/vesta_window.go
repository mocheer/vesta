//go:build window

// https://github.com/go-vgo/robotgo/issues/641
// @see https://github.com/robotn/gohook/issues/27
// 如果想要编译到linux，需要安装 mingw-w64-gcc,运行 CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ go build
package vesta

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/go-vgo/robotgo"
)

// 弃用
func (v *Vesta) AddEventMouseClickXY() *Vesta {

	action := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(ev any) {
			// 监听从网页发送的消息
			switch e := ev.(type) {
			case *runtime.EventConsoleAPICalled:

				args := e.Args
				var eventName string
				json.Unmarshal(args[0].Value, &eventName)

				switch eventName {
				case "$click":
					var x, y float64
					json.Unmarshal(args[1].Value, &x)
					json.Unmarshal(args[2].Value, &y)
					go func() {
						log.Println("click", x, y)
						robotgo.Move(int(x)+1920+rand.Intn(10), int(y)+140+rand.Intn(2))
						robotgo.MoveSmooth(int(x)+1920+rand.Intn(10), int(y)+140+rand.Intn(2))
						//
						chromedp.Run(ctx,
							chromedp.MouseClickXY(x, y),
						)
					}()

				}
			}
		})
		return nil
	})
	return v.AddTask(action)
}
