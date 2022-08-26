package vesta

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/mocheer/pluto/pkg/ts/img"
)

// Vesta 不支持单元测试环境
type Vesta struct {
	ctx     context.Context
	cancel  context.CancelFunc
	actions []chromedp.Action
}

func New() *Vesta {
	ctx, _ := chromedp.NewContext(context.Background())
	// defer cancel()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	//
	return &Vesta{
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddTask
func (v *Vesta) AddTask(action chromedp.Action) *Vesta {
	v.actions = append(v.actions, action)
	return v
}

// Viewport 设置页面宽高，用于截图等
func (v *Vesta) Viewport(width, height int64) *Vesta {
	return v.AddTask(chromedp.EmulateViewport(width, height))
}

// Nav
func (v *Vesta) Nav(url string) *Vesta {
	return v.AddTask(chromedp.Navigate(url))
}

// Location
func (v *Vesta) Location(url string) *Vesta {
	return v.AddTask(chromedp.Location(&url))
}

// WaitVisible
func (v *Vesta) Wait(d time.Duration) *Vesta {
	return v.AddTask(chromedp.Sleep(d))
}

// WaitVisible
func (v *Vesta) WaitReady(sel any) *Vesta {
	return v.AddTask(chromedp.WaitReady(sel))
}

// WaitVisible
func (v *Vesta) WaitID(id string) *Vesta {
	return v.AddTask(chromedp.WaitVisible(id, chromedp.ByID))
}

// WaitQuery
func (v *Vesta) WaitQuery(id string) *Vesta {
	return v.AddTask(chromedp.WaitVisible(id, chromedp.ByQuery))
}

// Eval
func (v *Vesta) Eval(jsScript string, res any) *Vesta {
	// res 会获取最后一个表达式的值，可以用分号，逗号正常的编写复杂的js脚本，只要最后一个表达式是最后要获取的值就可以了。
	action := chromedp.Evaluate(jsScript, res, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
		// 支持promise
		return p.WithAwaitPromise(true)
	})
	return v.AddTask(action)
}

// EvalModule 支持CommonJS的模块化支持，获取模块抛出的对象
func (v *Vesta) EvalModule(jsScript string, res any) *Vesta {
	jsScript = fmt.Sprintf(`const module={};%s;module.exports`, jsScript)
	return v.Eval(jsScript, res)
}

// Screen  New().Viewport(1920,1080).Screen()
func (v *Vesta) Screen(res *[]byte) *Vesta {
	return v.AddTask(chromedp.CaptureScreenshot(res))
}

// Run
func (v *Vesta) Run() error {
	return chromedp.Run(v.ctx, v.actions...)
}

// @see https://github.com/chromedp/chromedp/issues/592
// C:\Users\Administrator\AppData\Local\Temp
// Cancel 取消, window/temp下生成chromedp-runner文件有时候不会被移除，日积月累容易导致硬盘空间不足 => 待验证
func (v *Vesta) Cancel() {
	chromedp.Cancel(v.ctx)
	v.cancel()
}

// GetValue
func (v *Vesta) Get(jsScript string, res any) error {
	err := v.Eval(jsScript, &res).Run()
	return err
}

// GetValue
func (v *Vesta) GetValue(jsScript string) (any, error) {
	var res any
	err := v.Eval(jsScript, &res).Run()
	return res, err
}

// GetScreen New().Viewport(1920,1080).Screen()
func (v *Vesta) GetScreen() ([]byte, error) {
	var res []byte
	v.Screen(&res)
	err := v.Run()
	return res, err
}

// Screen New().Viewport(1920,1080).Screen()
func (v *Vesta) GetImage() (*img.Img, error) {
	bs, err := v.GetScreen()
	if err == nil {
		return img.FromBytes(bs)
	}
	return nil, err
}
