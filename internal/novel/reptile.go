package novel

import (
	"github.com/mocheer/vesta/pkg/vesta"
)

// 妙书斋
func Reptile() {
	vm = vesta.New().Nav("http://wap.miaoshuzhai.net/xs/178294")

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

	doReptile(novels)
}

func doReptile(data []*NovelChapter) {

}
