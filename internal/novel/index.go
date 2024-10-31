package zcode

import (
	"fmt"

	"github.com/mocheer/pluto/pkg/ts/clock"
	"github.com/mocheer/vesta/pkg/localdb"
	"github.com/mocheer/vesta/pkg/vesta"
)

type NovelChapter struct {
	Name  string
	Level int
	Link  string `gorm:"-"`
}

var vm *vesta.Vesta

// TableName 设置表名
func (NovelChapter) TableName() string {
	return "studio.dmap_novel"
}

var db, _ = localdb.Open("data.db")

func Reptile() {
	vm = vesta.NewWithDefault().Nav("http://www.miaoshuzhai.net/xs/142242/54129321.html")

	novels := []*NovelChapter{}
	err := vm.Get(`
({title:document.querySelector('.bookname>h1').textContent,content:Array.from(document.querySelector("#content").children)
  .map((e) => {
    if(e.tagName == 'BR'){
      return '\n'
    }
    let content = window.getComputedStyle(e, "before").content;
    content = content.slice(1, content.length - 1);
    return content;
  })
  .join("");})
	`,
		&novels)
	//
	vm.Cancel()
	if err != nil {
		panic(err)
	}

	doReptile(novels, 1)
}

func doReptile(data []*NovelChapter, level int) {

	//
	for _, item := range data {
		item.Level = level
		if item.Link != "" {
			nextData := []*NovelChapter{}
		start:
			fmt.Println("开始抓取", item)
			vm := vesta.NewWithDefault()
			err := vm.Nav(item.Link).Get(`Array.from(document.querySelectorAll(".villagetr,.towntr,.countytr,.citytr")).map(e => {
				let [td1,td2,td3] = e.querySelectorAll('td');
				let code = td1.innerText
				let name = td3?td3.innerText:td2.innerText
				name = name.replace("\n", " ")
				let aTag = td1.querySelector('a')
				let link = ""
				if (aTag) {
					link = aTag.href
				}
				return { code , name, link }
			});`, &nextData)
			vm.Cancel()
			if err != nil {
				fmt.Println(err)
				goto start
			} else {
				fmt.Println("抓取结束", item)
				doReptile(nextData, level+1)
			}
		}
	}

	tx := db.Create(data)
	fmt.Println("插入成功", tx.RowsAffected, clock.Now().Val(), len(data))

}
