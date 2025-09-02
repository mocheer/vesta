package main

import (
	"log"
	"math"

	"github.com/mocheer/vesta/pkg/vesta/bd/coord"
)

func main() {

	// // 批量转换示例
	// points := []coord.Point{
	// 	{116.98669500871492, 36.64824694844293}, // 天安门WGS84
	// 	{116.92753198793665, 36.61304331888675}, // 东方明珠WGS84
	// }

	// 第一步：转BD09

	p := coord.DB09MctoBD09(&coord.Point{
		13024401.82,
		4365353.30,
	})

	// 第二步：转GCJ02
	p2 := coord.BD09toGCJ02(&p)
	log.Println(p2)
	// 第三步：转WGS84
	wgsLng, wgsLat := gcj02ToWGS84(p2.Lng, p2.Lat)
	log.Println(wgsLng, wgsLat)
}

// GCJ02转WGS84
func gcj02ToWGS84(lng, lat float64) (float64, float64) {
	a := 6378245.0
	ee := 0.00669342162296594323

	transformLat := func(x, y float64) float64 {
		ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
		ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
		ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
		ret += (160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
		return ret
	}

	transformLng := func(x, y float64) float64 {
		ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
		ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
		ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
		ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
		return ret
	}

	dlat := transformLat(lng-105.0, lat-35.0)
	dlng := transformLng(lng-105.0, lat-35.0)
	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * math.Pi)
	dlng = (dlng * 180.0) / (a / sqrtmagic * math.Cos(radlat) * math.Pi)
	mglat := lat + dlat
	mglng := lng + dlng
	return lng*2 - mglng, lat*2 - mglat
}
