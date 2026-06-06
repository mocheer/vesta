package main

import (
	"time"

	"github.com/mocheer/vesta/pkg/vesta"
)

func main() {
	reptile()
	select {}
}

func reptile() {
	vm := vesta.New().Edge().Head()
	vm.SaveAllImages()
	vm.Nav("https://www.aigei.com/view/74226.html")
	vm.Sleep(12 * time.Hour)
	vm.Run()
}
