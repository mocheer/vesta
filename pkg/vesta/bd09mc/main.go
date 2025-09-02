package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/pluto/pkg/ds/ds_json"
	"github.com/mocheer/xena/pkg/gm"
)

func main() {
	data := map[string]gm.Point{}
	ds.EachFiles("../baidu_pano/data", func(filename string, fi os.FileInfo) {
		if strings.HasSuffix(filename, "config.json") {
			var c Config
			ds_json.ReadFile(filename, &c)
			for _, content := range c.Content {
				lon, lat := toLonLat(float64(content.X)/100, float64(content.Y)/100)
				p := gm.Point{lon, lat}
				data[content.ID] = p
				// log.Println(content.ID, lon, lat)
			}
		}
	})
	ds_json.Save("coordinates.json", data)
}

func toLonLat(mcX, mcY float64) (float64, float64) {

	// 第一步：转BD09
	bdLng, bdLat, err := bd09mcToBD09(mcX, mcY)
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0, 0
	}

	// 第二步：转GCJ02
	gcjLng, gcjLat := bd09ToGCJ02(bdLng, bdLat)

	// 第三步：转WGS84
	wgsLng, wgsLat := gcj02ToWGS84(gcjLng, gcjLat)
	return wgsLng, wgsLat
}
