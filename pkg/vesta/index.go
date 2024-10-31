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

// SaveAllResource 用于保存系统加载的所有资源文件
func SaveAllResource(url string) {
	// SaveAllResource必须在Nav之前，否则最开始的html和js不会下载
	v := New().Head().SaveAllResource().Nav(url)
	defer v.Cancel()
	v.Sleep(time.Second * 3000)
	v.Run()
}
