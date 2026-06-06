package zcode

import (
	"fmt"

	"github.com/mocheer/pluto/pkg/ts/clock"
	"github.com/mocheer/vesta/pkg/localdb"
	"github.com/mocheer/vesta/pkg/vesta"
)

type ZCode struct {
	Code  string `gorm:"column:id;primary_key"`
	Name  string
	Level int
	Link  string `gorm:"-"`
}

var vm *vesta.Vesta

// TableName 设置表名
func (ZCode) TableName() string {
	return "studio.dmap_zcode"
}

var db, _ = localdb.Open("data.db")

func Reptile() {
	// https://data.stats.gov.cn/dg/website/page.html#/pc/national/countryYearData
	// https://www.stats.gov.cn/xxgk/tjbz/gjtjbz/201310/P020200612582963846218.PDF
	vm = vesta.New().Nav("http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2021/index.html")

	provinces := []*ZCode{}
	err := vm.Get(`
		Array.from(document.querySelectorAll("body table table a")).map(e=>({name:e.innerText.replace("\n"," "),link:e.href,code:e.href.split('/').pop().slice(0,2)}))
	`,
		&provinces)
	//
	vm.Cancel()
	if err != nil {
		panic(err)
	}

	doReptile(provinces, 1)
}

func doReptile(data []*ZCode, level int) {

	//
	for _, item := range data {
		item.Level = level
		if item.Link != "" {
			nextData := []*ZCode{}
		start:
			fmt.Println("开始抓取", item)
			vm := vesta.New()
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
