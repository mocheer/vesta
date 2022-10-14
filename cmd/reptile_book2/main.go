package main

import (
	"os"
	"sort"
	"strings"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/samber/lo"
)

func main() {
	urls := []string{}
	ds.EachFiles("test", func(filename string, fi os.FileInfo) {
		urls = append(urls, filename)
	})
	sort.Strings(urls)
	texts := lo.Map(urls, func(url string, _ int) string {
		data, _ := ds.ReadFile(url)
		return string(data)
	})
	ds.Save("book.txt", []byte(strings.Join(texts, "\n====\n")))
}
