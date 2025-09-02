package main

import (
	"fmt"
	"math"
)

const (
	earthRadius = 6378137.0             // 地球半径（米）
	maxMercator = earthRadius * math.Pi // 墨卡托投影最大坐标值（赤道周长的一半）
	tileSize    = 256.0                 // 瓦片大小（像素）
)

// 将BD09经纬度转换为墨卡托投影坐标（米）
func lonLatToMercator(lng, lat float64) (x, y float64) {
	lngRad := lng * math.Pi / 180.0
	latRad := lat * math.Pi / 180.0

	x = earthRadius * lngRad
	y = earthRadius * math.Log(math.Tan(math.Pi/4+latRad/2))
	return
}

// 计算像素向量（dx, dy）
func GetPixelVector(lng1, lat1, lng2, lat2 float64, zoom int) (float64, float64) {
	// 经纬度转墨卡托坐标
	x1, y1 := lonLatToMercator(lng1, lat1)
	x2, y2 := lonLatToMercator(lng2, lat2)

	// 计算归一化坐标（0-1范围）
	normalizedX1 := (x1 + maxMercator) / (2 * maxMercator)
	normalizedY1 := (y1 + maxMercator) / (2 * maxMercator)
	normalizedX2 := (x2 + maxMercator) / (2 * maxMercator)
	normalizedY2 := (y2 + maxMercator) / (2 * maxMercator)

	// 计算当前缩放级别的总像素数
	mapSize := math.Pow(2, float64(zoom)) * tileSize

	// 计算像素坐标（Y轴需要翻转）
	pixelX1 := normalizedX1 * mapSize
	pixelY1 := (1 - normalizedY1) * mapSize
	pixelX2 := normalizedX2 * mapSize
	pixelY2 := (1 - normalizedY2) * mapSize

	// 返回像素向量
	return pixelX2 - pixelX1, pixelY2 - pixelY1
}

func main() {
	// 示例：计算北京到上海在缩放级别10下的像素向量
	beijing := [2]float64{180, 0} // 北京经纬度
	shanghai := [2]float64{0, 0}  // 上海经纬度
	zoom := 21                    // 缩放级别

	dx, dy := GetPixelVector(
		beijing[0], beijing[1],
		shanghai[0], shanghai[1],
		zoom,
	)

	fmt.Printf("像素向量: (%.2f, %.2f) 像素\n", dx, dy)
}
