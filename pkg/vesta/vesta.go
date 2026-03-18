package vesta

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/pluto/pkg/ts/ctp"
	"github.com/mocheer/pluto/pkg/ts/img"
)

type Vesta struct {
	options []chromedp.ExecAllocatorOption
	cancels []context.CancelFunc
	actions []chromedp.Action
	//
	ctx context.Context
}

// initConnect
func (v *Vesta) initConnect() {
	// 创建一个新的执行分配器（ExecAllocator），负责启动浏览器并管理其生命周期
	// 这个分配器包含了启动和配置 Chrome 浏览器所需的所有设置和资源。
	// TODO:每次都是重新启动一个新浏览器进程，没有复用之前的cookie，经常要重新登录
	// TODO:支持用 NewRemoteAllocator 来连接到一个已经在运行的 Chrome 实例，而不是启动一个新的浏览器实例。=> 重用同一个浏览器实例
	//      需要确保开启了远程调试功能。这通常通过在启动 Chrome 时使用 --remote-debugging-port 参数来实现。
	// chrome --remote-debugging-port=9222
	// allocator, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")

	// 根context
	ctx := context.Background()
	allocator, cancel := chromedp.NewExecAllocator(ctx, v.options...)
	v.cancels = append(v.cancels, cancel)
	// 浏览器实例的根context
	ctx, cancel = chromedp.NewContext(allocator)
	v.cancels = append(v.cancels, cancel)
	v.ctx = ctx
}

// connect
func (v *Vesta) connect() error {
	ctx := context.Background()
	// 尝试连接远程浏览器
	allocator, cancel := chromedp.NewRemoteAllocator(ctx, "http://localhost:9222")
	v.cancels = append(v.cancels, cancel)

	// 浏览器实例的根context
	ctx, cancel = chromedp.NewContext(allocator)
	v.cancels = append(v.cancels, cancel)
	// 测试连接，如果有错误，直接取消关闭
	if err := chromedp.Run(ctx, chromedp.Navigate("about:blank")); err != nil {
		v.Cancel()
		return err
	}
	v.ctx = ctx
	return nil
}

// initContext
func (v *Vesta) initContext() {
	if v.ctx == nil {
		err := v.connect()
		if err != nil {
			log.Println("无法连接远程浏览器，正在重新启动...")
			v.initConnect()
		}
	}
}

// Head
// 设置浏览器
func (v *Vesta) Head() *Vesta {
	v.options = append(v.options, chromedp.Flag("headless", false))
	return v
}

// Edge
// 使用Edge浏览器
func (v *Vesta) Edge() *Vesta {
	v.options = append(v.options, chromedp.ExecPath("C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe"))
	return v
}

// 这个是F11全屏，连浏览器工具栏都隐藏了
func (v *Vesta) Fullscreen() *Vesta {
	v.options = append(v.options, chromedp.Flag("start-fullscreen", true))
	return v
}

// Size
// 设置浏览器大小
// chromedp.Flag("window-size", "1900,1000")
func (v *Vesta) Size(width, height int) *Vesta {
	v.options = append(v.options, chromedp.WindowSize(width, height))
	return v
}

// ProxyServer
func (v *Vesta) ProxyServer(url string) *Vesta {
	v.options = append(v.options, chromedp.ProxyServer(url))
	return v
}

// UserAgent
func (v *Vesta) UserAgent(agent string) *Vesta {
	v.options = append(v.options, chromedp.UserAgent(agent))
	return v
}

// NewContext 用于在同一个浏览器打开新页签
func (v *Vesta) NewContext() *Vesta {
	v.initContext()
	cancels := []context.CancelFunc{}
	ctx, cancel := chromedp.NewContext(v.ctx)
	cancels = append(cancels, cancel)
	return &Vesta{
		ctx:     ctx,
		cancels: cancels,
	}
}

// GetAllocator
func (v *Vesta) GetAllocator() chromedp.Allocator {
	return chromedp.FromContext(v.ctx).Allocator
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
	v.cancels = []context.CancelFunc{}
}

// AddTask
func (v *Vesta) AddTask(action chromedp.Action) *Vesta {
	v.actions = append(v.actions, action)
	return v
}

// Reload
// 重新加载当前标签页
func (v *Vesta) Reload() *Vesta {
	return v.AddTask(chromedp.Reload())
}

// ClosePage
// 关闭当前标签页
func (v *Vesta) ClosePage() *Vesta {
	return v.AddTask(page.Close())
}

// Viewport 设置页面宽高，用于截图等
// 这个不是浏览器的宽高，仅仅只是网页内容的宽高
// 这个值一般小于浏览器宽高，大于没有意义
func (v *Vesta) Viewport(width, height int64) *Vesta {
	return v.AddTask(chromedp.EmulateViewport(width, height))
}

// Setheaders 可以在Nav访问之前设置请求头
// 比如 v.Setheaders(map[string]any{"X-Header": "my request header"}).Nav('http://localhost') 这样在服务端可以获取到这个请求头的信息
func (v *Vesta) Setheaders(headers map[string]any) *Vesta {
	return v.AddTask(network.Enable()).AddTask(network.SetExtraHTTPHeaders(network.Headers(headers)))
}

// Nav
// 将当前页跳转至目标页面
// TODO 支持模拟Origin和Host
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
	return v.AddTask(
		chromedp.WaitVisible(
			id,
			chromedp.ByID,
		),
	)
}

// WaitQuery selector 为css选择器
func (v *Vesta) WaitQuery(selector string) *Vesta {
	return v.AddTask(chromedp.WaitVisible(selector, chromedp.ByQuery))
}

// Eval
// 这个是在网页的load事件后执行的，如果要提前执行，需要用 chromedp.ActionFunc(似乎无效)
// js对象可以直接用结构体接收数据
func (v *Vesta) Eval(jsScript string, res any) *Vesta {
	// res 会获取最后一个表达式的值，可以用分号，逗号正常的编写复杂的js脚本，只要最后一个表达式是最后要获取的值就可以了。
	action := chromedp.Evaluate(jsScript, res, GetDefaultEvaluateOptions())
	return v.AddTask(action)
}

// FetchJSON
func (v *Vesta) FetchJSON(url string, config string, res any) *Vesta {
	return v.Eval(fmt.Sprintf(`fetch("%s",%s).then(res=>res.json())`, url, config), res)
}

// chromedp.EvaluateAsDevTools 允许你执行JavaScript代码，就像在开发者工具控制台中执行一样
// 相比于Evaluate，这个功能能访问特定的DevTools的一些API，可以使用 DevTools 专属的快捷方法
// DevTools 上下文切换有额外开销，普通任务用 Evaluate 更高效。
// 例子：
// DevTools 的 $ 选择器：$('div') 相当于 document.querySelector， $$('div') 相当于 document.querySelectorAll
// DevTools 内部属性：可以返回非序列化的 DOM 节点引用
// DevTools 的 getEventListeners
// DevTools 上下文中，Error 对象有更完整的堆栈信息
// 获取元素引用，返回 DOM 节点引用，不是 JSON，res可以用any接收，然后跨越脚本传递 == 需要验证
// DevTools 性能指标
// DevTools 的 timeline API
func (v *Vesta) EvaluateAsDevTools(jsScript string, res any) *Vesta {
	// res 会获取最后一个表达式的值，可以用分号，逗号正常的编写复杂的js脚本，只要最后一个表达式是最后要获取的值就可以了。
	action := chromedp.EvaluateAsDevTools(jsScript, res, GetDefaultEvaluateOptions())
	return v.AddTask(action)
}

// Inject
// page.addScriptToEvaluateOnNewDocument 可以向页面注入脚本，这些脚本将在新文档创建时执行，即在页面的任何脚本执行之前，可以操作文档 document 对象，但html解析未开始
// 可用于重写原型链，如Array.prototype.push，重写xhr和fetch
// 可用于重写页面document.documentElement.innerHTML，但无法获取dom对象， document.body为空，但可以创建dom节点，只是只能通过 document.documentElement.appendChild 添加节点，跟body处于同一个级别，似乎可以添加script节点，因为浏览器会自动纠正，相当于在body中执行
// sketchfab会覆盖console.log，这里也可用于提前缓存再还原
// 这个在nav之前使用inject，
// 执行时机：
//
// 页面开始加载
// ↓
// 创建新的 Document 对象
// ↓
// ✅ `AddScriptToEvaluateOnNewDocument` 脚本执行 ← 关键时机！
// ↓
// 解析 HTML（包括 script 标签）
// ↓
// 按照文档顺序执行 script 标签（内联/外部）
// ↓
// # DOMContentLoaded 事件
// ↓
// 页面资源加载完成
// ↓
// load 事件
func (v *Vesta) Inject(jsScript string) *Vesta {
	action := chromedp.ActionFunc(func(cxt context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument(jsScript).Do(cxt)
		if err != nil {
			return err
		}
		return nil
	})
	return v.AddTask(action)
}

// TODO 验证
func (v *Vesta) InjectGuard() *Vesta {
	return v.Inject(`
Object.defineProperty(window, 'location', {
    get() {
        const proxy = new Proxy(originalLocation, {
            get(target, prop) {
                if (prop === 'href' || prop === 'assign' || prop === 'replace') {
                    return function() {
                        console.warn('页面跳转被阻止:', arguments[0]);
                        return undefined;
                    };
                }
                return target[prop];
            },
            set(target, prop, value) {
                if (prop === 'href') {
                    console.warn('location.href 设置被阻止:', value);
                    return true;
                }
                target[prop] = value;
                return true;
            }
        });
        return proxy;
    },
    set(value) {
        console.warn('window.location 赋值被阻止');
        return false;
    },
    configurable: false,
    enumerable: true
});

// 2. 拦截 window.open
window.open = function() {
    console.warn('window.open 被阻止:', arguments[0]);
    return null;
};

// 3. 拦截 form 提交
HTMLFormElement.prototype.submit = function() {
    console.warn('表单提交被阻止');
    return false;
};

// 4. 拦截 a 标签点击（可选）
document.addEventListener('click', function(e) {
    if (e.target.tagName === 'A' && e.target.href) {
        e.preventDefault();
        console.warn('链接点击被阻止:', e.target.href);
    }
}, true);
`)
}

type InterceptRequestParams struct {
	UpdateBody  map[string]string
	ReplaceBody map[string]string
	Match       string
}

// InterceptRequestWithJS 拦截请求
// TODO 将所有匹配的js都用iife的方式加载，并注入新的bom局部变量，来实现固有变量的重写
func (v *Vesta) InterceptRequestWithJS(args *InterceptRequestParams) *Vesta {
	action := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(e interface{}) {
			switch ev := e.(type) {
			// 需要设置fetch.Eable请求才会暂停
			case *fetch.EventRequestPaused:
				// 包括js、css等资源的请求都会进来
				// https://chromedevtools.github.io/devtools-protocol/tot/Fetch/#event-requestPaused
				go func(ee *fetch.EventRequestPaused) {
					request := ee.Request

					// 重写返回的结果
					updateScript, ok := args.UpdateBody[request.URL]

					if !ok && args.Match != "" {
						updateScript = args.Match
						ok = true
					}
					if ok {
						// ee.Request.Headers["[Accept"]
						UpdateFetchBody(ctx, ee, func(b []byte, r *http.Response) []byte {
							var body string
							code := fmt.Sprintf(`JSON.stringify((%s)(%s))`, updateScript, string(b))
							err := chromedp.Run(ctx, chromedp.Evaluate(code, &body))
							if err != nil {
								fmt.Println(err)
							}
							return []byte(body)
						})
						return
					}
					//
					body, ok := args.ReplaceBody[ee.Request.URL]
					if ok {
						ReplaceFetchBody(ctx, ee, body)
						return
					}
					c := chromedp.FromContext(ctx)
					ctxFetch := cdp.WithExecutor(ctx, c.Target)
					fetch.ContinueRequest(ee.RequestID).Do(ctxFetch)

				}(ev)
			}
		})
		return nil
	})
	// 启用发出requestPaused事件。请求将被暂停，直到客户端 调用failRequest、fulfilled request或continuerrequest /continueWithAuth中的一个。
	// 这里通过pattern指定拦截的url
	return v.AddTask(fetch.Enable()).AddTask(action)
}

// InterceptRequestWithBauduPano 拦截百度全景图
func (v *Vesta) InterceptRequestWithBauduPano(name string) *Vesta {

	action := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(e interface{}) {
			switch ev := e.(type) {
			// 需要设置fetch.Eable请求才会暂停
			case *fetch.EventRequestPaused:
				// 包括js、css等资源的请求都会进来
				// https://chromedevtools.github.io/devtools-protocol/tot/Fetch/#event-requestPaused
				go func(ee *fetch.EventRequestPaused) {
					request := ee.Request

					if strings.Contains(request.URL, "https://mapsv0.bdimg.com/?qt=sdata") {
						var group sync.WaitGroup
						u, _ := url.Parse(request.URL)
						query := u.Query()
						sid := query.Get("sid")
						//
						log.Println(sid)

						data := FetchByHttp(ee)
						config := string(data)
						configByte := []byte(config[strings.Index(config, "(")+1 : strings.LastIndex(config, ")")])
						ds.Save(fmt.Sprintf("./data/%s/%s/config.json", name, sid), configByte)

						auth := query.Get("auth")
						seckey := query.Get("seckey")
						udt := query.Get("udt")

						for i := range 4 {
							for j := range 8 {
								if i == 0 && j == 0 {
									pdataURL := fmt.Sprintf("https://mapsv1.bdimg.com/?qt=pdata&sid=%s&pos=%d_%d&z=%d&udt=%s&from=PC&auth=%s&seckey=%s", sid, i, j, 1, udt, auth, seckey)
									fname := fmt.Sprintf("./data/%s/%s/0_%d_%d.jpg", name, sid, i, j)
									if !ds.IsExist(fname) {
										group.Go(func() {
											log.Println(fname)
											ctp.Save(pdataURL, fname)
										})
									}
								}
								pdataURL := fmt.Sprintf("https://mapsv1.bdimg.com/?qt=pdata&sid=%s&pos=%d_%d&z=%d&udt=%s&from=PC&auth=%s&seckey=%s", sid, i, j, 4, udt, auth, seckey)
								filename := fmt.Sprintf("./%s/%s/%d_%d.jpg", name, sid, i, j)
								if !ds.IsExist(filename) {
									group.Go(func() {
										log.Println(filename)
										ctp.Save(pdataURL, filename)
									})

								}
							}
						}
						group.Wait()
					}
					c := chromedp.FromContext(ctx)
					ctxFetch := cdp.WithExecutor(ctx, c.Target)
					fetch.ContinueRequest(ee.RequestID).Do(ctxFetch)

				}(ev)
			}
		})
		return nil
	})
	// 启用发出requestPaused事件。请求将被暂停，直到客户端 调用failRequest、fulfilled request或continuerrequest /continueWithAuth中的一个。
	// 这里通过pattern指定拦截的url
	return v.AddTask(fetch.Enable()).AddTask(action)
}

func (v *Vesta) MouseClickXY(x, y float64) *Vesta {
	return v.AddTask(chromedp.MouseClickXY(x, y))
}

// SaveAllResource
func (v *Vesta) SaveAllResource() *Vesta {

	action := chromedp.ActionFunc(func(ctx context.Context) error {
		// 会有线程安全的问题
		urlMap := map[network.RequestID]string{}
		mu := &sync.Mutex{}
		chromedp.ListenTarget(ctx, func(e interface{}) {
			switch ev := e.(type) {
			case *network.EventRequestWillBeSent:
				mu.Lock()
				urlMap[ev.RequestID] = ev.Request.URL
				mu.Unlock()
			// 虽然开始接收Response，但可能数据未全部加载完
			// case *network.EventResponseReceived:
			// 	urlMap[ev.RequestID] = ev.Response.URL
			case *network.EventLoadingFinished:
				go func(ev *network.EventLoadingFinished, ctx context.Context) {
					mu.Lock()
					upath, ok := urlMap[ev.RequestID]
					mu.Unlock()
					if ok {
						c := chromedp.FromContext(ctx)
						ctxFetch := cdp.WithExecutor(ctx, c.Target)
						data, err := network.GetResponseBody(ev.RequestID).Do(ctxFetch)
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Println(upath)
						// 有可能是一个blob:开头的自定义url字符串
						// 这种自定义字符串其实可以直接过滤，因为每次都不一样
						// 有一种可能的需求是想要拦截获取最终的blob资源数据，这个时候要另外处理（原始数据可能会加密，直接查看blob数据更清晰）
						if strings.HasPrefix(upath, "blob:") {
							// path = strings.Replace(path, "blob:", "", 1)
							return
						}
						uri, err := url.Parse(upath)
						if err != nil {
							fmt.Println(err)
						}
						filename := "./testdata/save/" + uri.Hostname() + "/" + uri.Path

						// url接口，很多都是相同名称不同参数
						if len(uri.RawQuery) > 0 {
							// 但这样做会导致404
							filename += "@" + url.QueryEscape(uri.RawQuery)
						}
						if strings.HasSuffix(filename, "/") {
							filename += "index.html"
						}
						//
						err = ds.Save(filename, data)
						if err != nil {
							fmt.Println(err)
						}
					}
				}(ev, ctx)
			case *network.EventWebSocketCreated:
				log.Println("EventWebSocketCreated", ev)
			case *network.EventWebSocketClosed:
				log.Println("EventWebSocketClosed", ev)
			case *network.EventWebSocketWillSendHandshakeRequest:
				log.Println("EventWebSocketWillSendHandshakeRequest", ev)
			case *network.EventWebSocketHandshakeResponseReceived:
				log.Println("EventWebSocketHandshakeResponseReceived", ev)
			case *network.EventWebSocketFrameReceived:
				log.Println("EventWebSocketFrameReceived", ev.Response)
			case *network.EventWebSocketFrameSent:
				log.Println("EventWebSocketFrameSent", ev.Response)
			case *network.EventWebSocketFrameError:
				log.Println("EventWebSocketFrameError", ev)

			}
		})
		return nil
	})
	// 不能开启 fetch.Enable(),否则会被阻塞，需要监听暂停事件并重新请求并返回结果
	return v.AddTask(action)
}

// EvalModule 支持CommonJS的模块化支持，获取模块抛出的对象
func (v *Vesta) EvalModule(jsScript string, res any) *Vesta {
	jsScript = fmt.Sprintf(`const module={};%s;module.exports`, jsScript)

	return v.Eval(jsScript, res)
}

// Screen
// New().Viewport(1920,1080).Screen()
func (v *Vesta) Screen(res *[]byte) *Vesta {
	return v.AddTask(chromedp.CaptureScreenshot(res))
}

// Run
func (v *Vesta) Run() error {
	v.initContext()
	// 任务执行的超时时间：默认5分钟
	// 脚本可能是Promise等待超长，又或者远程服务器超慢，这个时候需要设置超时时间
	ctx, cancel := context.WithTimeout(v.ctx, 5*time.Minute)
	defer cancel()
	//
	err := chromedp.Run(ctx, v.actions...)
	v.actions = []chromedp.Action{}
	return err
}

// GetValue
// json 可以直接用结构体接收数据
func (v *Vesta) Get(jsScript string, res any) *Vesta {
	err := v.Eval(jsScript, &res)
	return err
}

// GetValue
func (v *Vesta) GetValue(jsScript string) (any, error) {
	var res any
	v.Eval(jsScript, &res)
	err := v.Run()
	return res, err
}

// GetScreen
// New().Screen()
func (v *Vesta) RunGetScreen() ([]byte, error) {
	var res []byte
	v.Screen(&res)
	err := v.Run()
	return res, err
}

// GetImage
// New().RunGetImage()
func (v *Vesta) RunGetImage() (*img.Img, error) {
	bs, err := v.RunGetScreen()
	if err == nil {
		i, _, err := img.FromBytes(bs)
		return i, err
	}
	return nil, err
}

func (v *Vesta) AddEventMouseClickXY() *Vesta {
	return v
}
