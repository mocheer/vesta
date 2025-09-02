package main

import (
	"fmt"
	"log"
	"math"
)

const (
	baiduPI       = math.Pi * 3000.0 / 180.0 // 百度坐标系专用PI
	a             = 6378245.0                // WGS84椭球长半轴
	ee            = 0.00669342162296594323   // WGS84椭球第一偏心率平方
	earthRadius   = 6378137.0                // 地球半径（米）
	mercatorRange = 20037508.34              // 墨卡托投影范围
)

// Point 定义经纬度点结构
type Point struct {
	Lng, Lat float64 // 经度, 纬度
}

// MercatorPoint 定义墨卡托投影坐标点
type MercatorPoint struct {
	X, Y float64 // 投影坐标（米）
}

// PixelPoint 定义像素坐标点
type PixelPoint struct {
	X, Y float64 // 像素坐标
}

// delta 计算WGS84到GCJ02的坐标偏移量
func delta(lng, lat float64) (float64, float64) {
	dLat := transformLat(lng-105.0, lat-35.0)
	dLng := transformLng(lng-105.0, lat-35.0)

	radLat := lat * math.Pi / 180.0
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)

	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * math.Pi)
	dLng = (dLng * 180.0) / (a / sqrtMagic * math.Cos(radLat) * math.Pi)

	return dLng, dLat
}

// transformLat 纬度变换辅助函数
func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y +
		0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) +
		20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) +
		40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) +
		320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

// transformLng 经度变换辅助函数
func transformLng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x +
		0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) +
		20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) +
		40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) +
		300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

// WGS84ToGCJ02 将WGS84坐标转换为GCJ02(火星坐标系)
func WGS84ToGCJ02(lng, lat float64) (float64, float64) {
	if lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271 {
		return lng, lat // 中国境外不转换
	}
	dLng, dLat := delta(lng, lat)
	return lng + dLng, lat + dLat
}

// GCJ02ToMercator 将GCJ02坐标转换为墨卡托投影坐标（米）
func GCJ02ToMercator(lng, lat float64) (float64, float64) {
	// 将经纬度转换为弧度
	radLng := lng * math.Pi / 180.0
	radLat := lat * math.Pi / 180.0

	// 墨卡托投影公式
	x := earthRadius * radLng
	y := earthRadius * math.Log(math.Tan(math.Pi/4.0+radLat/2.0))

	return x, y
}

// MercatorToPixel 将墨卡托坐标转换为像素坐标
func MercatorToPixel(x, y float64, zoom int) (float64, float64) {
	resolution := mercatorRange * 2 / math.Pow(2, float64(zoom)) / 256.0
	pixelX := (x + mercatorRange) / resolution
	pixelY := (mercatorRange - y) / resolution
	return pixelX, pixelY
}

// BD09转GCJ02
func bd09ToGCJ02(lng, lat float64) (float64, float64) {
	x := lng - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*math.Pi*3000.0/180.0)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*math.Pi*3000.0/180.0)
	gcjLng := z * math.Cos(theta)
	gcjLat := z * math.Sin(theta)
	return gcjLng, gcjLat
}

// GCJ02ToBD09 将GCJ02坐标转换为BD09(百度坐标系)
func GCJ02ToBD09(lng, lat float64) (float64, float64) {
	z := math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*baiduPI)
	theta := math.Atan2(lat, lng) + 0.000003*math.Cos(lng*baiduPI)
	bdLng := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006
	return bdLng, bdLat
}

// WGS84ToBaiduProjection 将WGS84坐标转换为百度投影坐标
func WGS84ToBaiduProjection(lng, lat float64, zoom int) (float64, float64) {
	// 转换到GCJ02坐标系
	// lng, lat = WGS84ToGCJ02(lng, lat)
	// lng, lat = GCJ02ToBD09(lng, lat)

	// 转换为墨卡托投影坐标
	mercatorX, mercatorY := GCJ02ToMercator(lng, lat)
	log.Println(mercatorX, mercatorY)
	// 转换为像素坐标
	pixelX, pixelY := MercatorToPixel(mercatorX, mercatorY, zoom)

	return pixelX, pixelY
}

// 批量转换函数
func BatchConvertToBaiduProjection(points []Point, zoom int) []PixelPoint {
	results := make([]PixelPoint, len(points))
	for i, p := range points {
		x, y := WGS84ToBaiduProjection(p.Lng, p.Lat, zoom)
		results[i] = PixelPoint{X: x, Y: y}
	}
	return results
}

// 计算指定zoom层级的地图分辨率（米/像素）
func ResolutionAtZoom(zoom int) float64 {
	return mercatorRange * 2 / math.Pow(2, float64(zoom)) / 256.0
}

func main() {
	// 示例坐标（北京天安门）
	beijing := Point{116.397428, 39.90923}

	// 不同zoom层级的转换
	zooms := []int{10, 15, 18}

	for _, zoom := range zooms {
		pixelX, pixelY := WGS84ToBaiduProjection(beijing.Lng, beijing.Lat, zoom)
		resolution := ResolutionAtZoom(zoom)

		fmt.Printf("Zoom %d (分辨率: %.2f 米/像素):\n", zoom, resolution)
		fmt.Printf("  像素坐标: (%.2f, %.2f)\n", pixelX, pixelY)
		fmt.Println()
	}

	// 批量转换示例
	points := []Point{
		{116.98670236090803, 36.648249044557204},
		{116.92753994214675, 36.61305052194565},
		{0, 0},
		{180, 0},
	}
	// 227525136
	// 227525136
	zoom := 21
	results := BatchConvertToBaiduProjection(points, zoom)

	fmt.Printf("批量转换结果 (Zoom %d):\n", zoom)
	for i, p := range results {
		fmt.Printf("坐标点 %d: (%.2f, %.2f)\n", i+1, p.X, p.Y)
	}
}
