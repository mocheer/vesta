package vesta

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/mocheer/pluto/ts/img"
)

// vesta 不支持单元测试环境
type vesta struct {
	ctx     context.Context
	cancel  context.CancelFunc
	actions []chromedp.Action
}

func New() *vesta {
	ctx, cancel := chromedp.NewContext(context.Background())
	return &vesta{
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddTask
func (v *vesta) AddTask(action chromedp.Action) *vesta {
	v.actions = append(v.actions, action)
	return v
}

// Viewport 设置页面宽高，用于截图等
func (v *vesta) Viewport(width, height int64) *vesta {
	return v.AddTask(chromedp.EmulateViewport(width, height))
}

// Nav
func (v *vesta) Nav(url string) *vesta {
	return v.AddTask(chromedp.Navigate(url))
}

// WaitVisible
func (v *vesta) Wait(d time.Duration) *vesta {
	return v.AddTask(chromedp.Sleep(d))
}

// WaitVisible
func (v *vesta) WaitID(id string) *vesta {
	return v.AddTask(chromedp.WaitVisible(id, chromedp.ByID))
}

// WaitQuery
func (v *vesta) WaitQuery(id string) *vesta {
	return v.AddTask(chromedp.WaitVisible(id, chromedp.ByQuery))
}

// Eval
func (v *vesta) Eval(jsScript string, res interface{}) *vesta {
	// res 会获取最后一个表达式的值，可以用分号，逗号正常的编写复杂的js脚本，只要最后一个表达式是最后要获取的值就可以了。
	action := chromedp.Evaluate(jsScript, res, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
		// 支持promise
		return p.WithAwaitPromise(true)
	})
	return v.AddTask(action)
}

// EvalModule 支持CommonJS的模块化支持，获取模块抛出的对象
func (v *vesta) EvalModule(jsScript string, res interface{}) *vesta {
	jsScript = fmt.Sprintf(`const module={};%s;module.exports`, jsScript)
	return v.Eval(jsScript, res)
}

// Screen  New().Viewport(1920,1080).Screen()
func (v *vesta) Screen(res *[]byte) *vesta {
	return v.AddTask(chromedp.CaptureScreenshot(res))
}

// Run
func (v *vesta) Run() error {
	defer v.cancel()
	return chromedp.Run(v.ctx, v.actions...)
}

// GetValue
func (v *vesta) GetValue(jsScript string) interface{} {
	var res interface{}
	err := v.Eval(jsScript, &res).Run()
	if err != nil {
		panic(err)
	}
	return res
}

// GetScreen New().Viewport(1920,1080).Screen()
func (v *vesta) GetScreen() []byte {
	var res []byte
	v.Screen(&res)
	err := v.Run()
	if err != nil {
		panic(err)
	}
	return res
}

// Screen New().Viewport(1920,1080).Screen()
func (v *vesta) GetImage() (*img.Img, error) {
	return img.FromBytes(v.GetScreen())
}
