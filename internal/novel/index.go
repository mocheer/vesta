package novel

import (
	"github.com/mocheer/vesta/pkg/localdb"
	"github.com/mocheer/vesta/pkg/vesta"
)

type NovelChapter struct {
	Name string
	Link string
}

var vm *vesta.Vesta

// TableName 设置表名
func (NovelChapter) TableName() string {
	return "studio.dmap_novel"
}

var db, _ = localdb.Open("data.db")
