package vesta

import "time"

func New() *Vesta {
	return &Vesta{
		options: GetDefaultExecAllocatorOptions(),
	}
}

// Screenshot 用于直接截图
func Screenshot(url string) ([]byte, error) {
	v := New().Nav(url)
	defer v.Cancel()
	return v.RunGetScreen()
}

// 支持脚本截图，根据脚本重新定位到最终的页面再截图
// TODO 支持多页面截图
func ScreenshotWithScript(url string, jsScript string) ([]byte, error) {
	v := New().Nav(url)
	defer v.Cancel()
	return v.RunGetScreen()
}

// SaveAllResource 用于保存系统加载的所有资源文件
func SaveAllResource(url string, dir string) {
	// SaveAllResource必须在Nav之前，否则最开始的html和js不会下载
	v := New().Head().SaveAllResource(dir).Nav(url)
	defer v.Cancel()
	v.Sleep(time.Second * 3000)
	v.Run()
}

// Script
// 用于支持长脚本,跨页面
// TODO 暂写伪代码
func Script(url, jsScript string) ([]byte, error) {
	v := New().Nav(url)

	v.GetValue(jsScript)

	type ScriptAction struct {
		URL    string
		Script string
	}
	type VestaExport struct {
		Actions []ScriptAction
	}
	v.Get("$vestaExport", &VestaExport{})

	defer v.Cancel()
	return nil, nil
}
