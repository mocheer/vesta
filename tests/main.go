package main

import (
	"fmt"

	"github.com/mocheer/vesta/cmd/vesta"
)

func main() {
	fmt.Println(vesta.New().GetValue("1+1"))
	fmt.Println(vesta.New().GetValue("'hello world '+ 'vesta'"))
	fmt.Println(vesta.New().GetValue("true"))
	fmt.Println(vesta.New().GetValue("({a:1})"))
	fmt.Println(vesta.New().Nav("https://www.baidu.com").GetValue("document.title"))
}
