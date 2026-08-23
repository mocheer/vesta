package main

import (
	"log"
	"time"

	"github.com/mocheer/pluto/pkg/ts/clock"
	"github.com/mocheer/vesta/pkg/vesta"
)

func main() {

	// reptile("./data/各种地形地面", "https://www.aigei.com/view/70659.html#items")//>65px
	// reptile("./data/地宫", "https://www.aigei.com/view/70677.html#items")
	// reptile("./data/高清建筑", "https://www.aigei.com/view/8665.html#items")
	// reptile("./data/仙境幻想", "https://www.aigei.com/view/8610.html?items") //缺少
	// reptile("./data/仙境幻想-场景素材", "https://www.aigei.com/view/8658.htm")
	reptile("./data/梦幻西游", "https://www.aigei.com/set/menghuanxiyouquantao-34468438.html")

	select {}
}

func reptile(dir string, url string) {
	vm := vesta.New().Edge().Head()
	vm.SaveAllImages(dir)
	vm.Nav(url)
	vm.Sleep(10 * time.Second)
	vm.Run()
	// 每秒点击一次
	clock.SetInterval(func() {
		log.Println("点击下一页")
		// // 每次都要重新设置这个task?这里点击下一页会丢失上下文
		// vm.SaveAllImages(dir)

		vm.Eval(`
		Array.from(document.querySelectorAll('ul.unit-content-main')).map(e=>e.querySelector('img')).filter(e=>e.src.includes('thumbnail')).forEach(e=>{
			e.click();
		});
		`, nil)
		// fp-prev 上一页
		// fp-next 下一页
		err := vm.Eval(`setTimeout(()=>{document.querySelector('.fp-next').click()},3000)`, nil).RunWithoutCtx()
		// 处理错误
		if err != nil {
			log.Println(err)
		}
	}, time.Second*8, false)

	select {}

}
