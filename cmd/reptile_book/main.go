package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/vesta/cmd/vesta"
)

func main() {

	doReptile := func(url string) string {
		vm := vesta.New()
		time.Sleep(time.Millisecond * 100)
		vm.Nav(url)

		var text string
		vm.Get(`
		document.getElementById('content').innerText
		`,
			&text)
		text = strings.Replace(text, "章节错误,点此报送(免注册), 报送后维护人员会在两分钟内校正章节内容,请耐心等待", "", -1)
		if text == "" {

			return url
		}
		var next string
		vm.Get(`
			document.getElementsByClassName('next')[0].href
		`,
			&next)

		ds.Save(filepath.Base(url)+".txt", []byte(text))

		return next
	}
	next := "https://www.kanshu5.net/172/172822/88966990.html"
	for ; next != ""; next = doReptile(next) {
		fmt.Println(next)
	}
	//

}
