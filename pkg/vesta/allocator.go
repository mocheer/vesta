package vesta

import (
	"github.com/chromedp/chromedp"
)

// 分配器选项
func GetDefaultExecAllocatorOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),      // 禁用网络安全，使得可以跨域操作iframe
		chromedp.Flag("ignore-certificate-errors", true), // 忽略错误，这样可以直接跳过https的认证错误页面)
		chromedp.NoSandbox, // 禁用沙盒，提升性能，可执行恶意代码
		//  chromedp.WithLogf(log.Printf)
		// chromedp.UserAgent(userAgent),                 // 有些网站是根据UserAgent进行拦截和路由跳转的
		// chromedp.NoFirstRun,                           // chrome首次运行常有欢迎界面和初始化设置，这里禁用
		// chromedp.ProxyServer(url),                     // 代理服务器
	)
}
