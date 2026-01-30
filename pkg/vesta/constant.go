package vesta

import (
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// 分配器选项
func GetDefaultExecAllocatorOptions() []chromedp.ExecAllocatorOption {
	cwd, _ := os.Getwd()
	userDataDir := filepath.Join(cwd, "chrome-data") //考虑采用浏览器的默认位置，
	//
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.IgnoreCertErrors,                            // 忽略错误，这样可以直接跳过https的认证错误页面)
		chromedp.Flag("enable-automation", false),            // 禁用自动化提示
		chromedp.Flag("disable-web-security", true),          // 禁用网络安全，使得可以跨域操作iframe
		chromedp.NoDefaultBrowserCheck,                       // 跳过默认浏览器检查
		chromedp.UserDataDir(userDataDir),                    // 设置用户数据目录，包含了用户的所有个性化设置和浏览数据，相当于浏览器的"用户配置文件"，包含cookie等，用于持久化
		chromedp.NoSandbox,                                   // 禁用沙盒，提升性能，可执行恶意代码
		chromedp.NoFirstRun,                                  // chrome首次运行常有欢迎界面和初始化设置，这里禁用
		chromedp.Flag("remote-debugging-port", "9222"),       // 开启调试端口
		chromedp.Flag("remote-debugging-address", "0.0.0.0"), // 允许远程连接
		// chromedp.ExecPath("C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe"), //  如果找不到会去默认的几个位置查找chrome.exe , 这里默认使用edge
		// chromedp.WithLogf(log.Printf)
		// chromedp.UserAgent(userAgent),                 // 有些网站是根据UserAgent进行拦截和路由跳转的
		// chromedp.ProxyServer(url),                     // 代理服务器
	)
}

func GetDefaultUserDir() string {
	// 使用Windows系统上Chrome的User Data的默认路径
	homeDir, _ := os.UserHomeDir() // 获取用户主目录,在windows是C:\Users\Administrator
	userDataDir := filepath.Join(homeDir, "AppData", "Local", "Google", "Chrome", "User Data")
	return userDataDir
}

// chromedp.Flag("ignore-certificate-errors", true), // 忽略错误，这样可以直接跳过https的认证错误页面)

func GetDefaultEvaluateOptions() chromedp.EvaluateOption {
	return func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
		// 返回值
		p.WithReturnByValue(true) //默认为true
		// 超时时间，默认不超时，这里要支持promise，所以设置超时时间3分钟
		// p.WithTimeout(runtime.TimeDelta(1000 * 60 * 3))
		// 支持promise
		return p.WithAwaitPromise(true)
	}
}
