package main

import (
	"fmt"

	"github.com/mocheer/vesta"
)

func main() {
	fmt.Println(vesta.New().GetBytes("1+1"))
	fmt.Println(vesta.New().GetString("'hello world '+ 'vesta'"))
	fmt.Println(string(vesta.New().GetBytes("({a:1})")))
	fmt.Println(vesta.New().Nav("https://www.baidu.com").GetString("document.title"))
	//
	fmt.Println(vesta.New().GetModuleString("module.exports = true"))
}
