package main

import (
	"fmt"
	"time"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/vesta/pkg/vesta"
)

func main() {
	v := vesta.NewWithHead().Nav("http://flood.control.command.dep003.devops.sc/capture-page")
	defer v.Cancel()
	v.Sleep(3500 * time.Millisecond)
	v.Eval("window.alert = null", nil)
	data, err := v.GetScreen()

	if err != nil {
		fmt.Println(err)
	}
	ds.Save("a.png", data)
}
