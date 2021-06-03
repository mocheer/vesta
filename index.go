package vesta

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/mocheer/pluto/img"
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
	// 好像是支持Promise的，待测试
	return v.AddTask(chromedp.Evaluate(jsScript, res))
}

// EvalModule 支持CommonJS的模块化支持，获取模块抛出的对象
func (v *vesta) EvalModule(jsScript string, res interface{}) *vesta {
	jsScript = fmt.Sprintf(`const module={};%s;module.exports`, jsScript)
	return v.Eval(jsScript, res)
}

// EvalModuleString 支持CommonJS的模块化支持，获取模块抛出的对象并强制转成string类型
func (v *vesta) EvalModuleString(jsScript string, res interface{}) *vesta {
	jsScript = fmt.Sprintf(`const module={};%s;String(module.exports)`, jsScript)
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

//  GetBytes
func (v *vesta) GetBytes(jsScript string) []byte {
	var res []byte
	v.Eval(jsScript, &res)
	err := v.Run()
	if err != nil {
		panic(err)
	}
	return res
}

//  GetString
func (v *vesta) GetString(jsScript string) string {
	var res string //这里不用byte，因为一些api的值是 unicode 编码，比如说document.title
	v.Eval(jsScript, &res)
	err := v.Run()
	if err != nil {
		panic(err)
	}
	return res
}

// GetModuleString
func (v *vesta) GetModuleString(jsScript string) string {
	var res string // 这里不用byte，因为一些api的值是 unicode 编码，比如说document.title
	v.EvalModuleString(jsScript, &res)
	err := v.Run()
	if err != nil {
		panic(err)
	}
	return res
}

// Screen New().Viewport(1920,1080).Screen()
func (v *vesta) GetImage() (*img.Picture, error) {
	var res []byte
	v.Screen(&res)
	err := v.Run()
	if err != nil {
		panic(err)
	}
	return img.FromBytes(res)
}
