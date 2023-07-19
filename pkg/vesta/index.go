package vesta

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/mocheer/pluto/pkg/ts/img"
)

// Vesta 不支持单元测试环境
type Vesta struct {
	ctx     context.Context
	cancels []context.CancelFunc
	actions []chromedp.Action
}

func New() *Vesta {
	cancels := []context.CancelFunc{}
	// 创建带有自定义标志选项的 Chromedp 实例
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),      //禁用网络安全，使得可以跨域操作iframe
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		// chromedp.Flag("headless", false),                 //开启图像界面
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	cancels = append(cancels, cancel)
	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background())
	cancels = append(cancels, cancel)
	return &Vesta{
		ctx:     ctx,
		cancels: cancels,
	}
}

func NewWithHead() *Vesta {
	cancels := []context.CancelFunc{}
	// 创建带有自定义标志选项的 Chromedp 实例
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),      //禁用网络安全，使得可以跨域操作iframe
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		chromedp.Flag("headless", false),                 //开启图像界面
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	cancels = append(cancels, cancel)
	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background())
	cancels = append(cancels, cancel)
	return &Vesta{
		ctx:     ctx,
		cancels: cancels,
	}
}

func NewWithProxy(url string) *Vesta {
	cancels := []context.CancelFunc{}
	// 创建带有自定义标志选项的 Chromedp 实例
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),      //禁用网络安全，使得可以跨域操作iframe
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		chromedp.Flag("headless", false),                 //开启图像界面
		chromedp.ProxyServer(url),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	cancels = append(cancels, cancel)
	ctx, cancel := chromedp.NewContext(allocCtx)
	// ctx, cancel := chromedp.NewContext(context.Background())
	cancels = append(cancels, cancel)
	return &Vesta{
		ctx:     ctx,
		cancels: cancels,
	}

}

// @see https://github.com/chromedp/chromedp/issues/592
// C:\Users\Administrator\AppData\Local\Temp
// Cancel 取消, window/temp下生成chromedp-runner文件有时候不会被移除，日积月累容易导致硬盘空间不足 => 待验证
// chromedp.Cancel(v.ctx)
// v.cancel()
func (v *Vesta) Cancel() {
	for i := len(v.cancels) - 1; i >= 0; i-- {
		v.cancels[i]()
	}
}

// WithTimeout 脚本可能是Promise等待超长，又或者远程服务器超慢，这个时候需要设置超时时间，注意这里是当前整个chromedp的超时时间，不是单个Nav的超时时间
func (v *Vesta) WithTimeout(timeout time.Duration) *Vesta {
	ctx, cancel := context.WithTimeout(v.ctx, timeout)
	v.ctx = ctx
	v.cancels = append(v.cancels, cancel)
	return v
}

// AddTask
func (v *Vesta) AddTask(action chromedp.Action) *Vesta {
	v.actions = append(v.actions, action)
	return v
}

// Reload
func (v *Vesta) Reload() *Vesta {
	return v.AddTask(chromedp.Reload())
}

// Viewport 设置页面宽高，用于截图等
func (v *Vesta) Viewport(width, height int64) *Vesta {
	return v.AddTask(chromedp.EmulateViewport(width, height))
}

// Setheaders 可以在Nav访问之前设置请求头
// 比如 v.Setheaders(map[string]any{"X-Header": "my request header"}).Nav('http://localhost') 这样在服务端可以获取到这个请求头的信息
func (v *Vesta) Setheaders(headers map[string]any) *Vesta {
	return v.AddTask(network.Enable()).AddTask(network.SetExtraHTTPHeaders(network.Headers(headers)))
}

// Nav
func (v *Vesta) Nav(url string) *Vesta {
	return v.AddTask(chromedp.Navigate(url))
}

// Location
func (v *Vesta) Location(url string) *Vesta {
	return v.AddTask(chromedp.Location(&url))
}

// Sleep
func (v *Vesta) Sleep(d time.Duration) *Vesta {
	return v.AddTask(chromedp.Sleep(d))
}

// WaitReady
func (v *Vesta) WaitReady(sel any) *Vesta {
	return v.AddTask(chromedp.WaitReady(sel))
}

// WaitID
func (v *Vesta) WaitID(id string) *Vesta {
	return v.AddTask(chromedp.WaitVisible(id, chromedp.ByID))
}

// WaitQuery selector 为css选择器
func (v *Vesta) WaitQuery(selector string) *Vesta {
	return v.AddTask(chromedp.WaitVisible(selector, chromedp.ByQuery))
}

// Eval
func (v *Vesta) Eval(jsScript string, res any) *Vesta {
	// res 会获取最后一个表达式的值，可以用分号，逗号正常的编写复杂的js脚本，只要最后一个表达式是最后要获取的值就可以了。
	action := chromedp.Evaluate(jsScript, res, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
		// p.WithReturnByValue(true)
		// p.WithTimeout()
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
	err := chromedp.Run(v.ctx, v.actions...)
	v.actions = []chromedp.Action{}
	return err
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
