package main

import (
	"fmt"
	"math"
)

// @see 百度安卓sdk中的Projection源码
// @see https://github.com/xyzj/toolbox/blob/main/coord/baidu.go#L43

// BD09MC转BD09 百度投影坐标系转bd09地理坐标系
func bd09mcToBD09(x, y float64) (float64, float64, error) {
	mcBands := []float64{12890594.86, 8362377.87, 5591021, 3481989.83, 1678043.12, 0}
	mc2ll := [][]float64{
		{1.410526172116255e-8, 0.00000898305509648872, -1.9939833816331, 200.9824383106796, -187.2403703815547, 91.6087516669843, -23.38765649603339, 2.57121317296198, -0.03801003308653, 17337981.2},
		{-7.435856389565537e-9, 0.000008983055097726239, -0.78625201886289, 96.32687599759846, -1.85204757529826, -59.36935905485877, 47.40033549296737, -16.50741931063887, 2.28786674699375, 10260144.86},
		{-3.030883460898826e-8, 0.00000898305509983578, 0.30071316287616, 59.74293618442277, 7.357984074871, -25.38371002664745, 13.45380521110908, -3.29883767235584, 0.32710905363475, 6856817.37},
		{-1.981981304930552e-8, 0.000008983055099779535, 0.03278182852591, 40.31678527705744, 0.65659298677277, -4.44255534477492, 0.85341911805263, 0.12923347998204, -0.04625736007561, 4482777.06},
		{3.09191371068437e-9, 0.000008983055096812155, 0.00006995724062, 23.10934304144901, -0.00023663490511, -0.6321817810242, -0.00663494467273, 0.03430082397953, -0.00466043876332, 2555164.4},
		{2.890871144776878e-9, 0.000008983055095805407, -3.068298e-8, 7.47137025468032, -0.00000353937994, -0.02145144861037, -0.00001234426596, 0.00010322952773, -0.00000323890364, 826088.5},
	}

	xAbs := math.Abs(x)
	yAbs := math.Abs(y)

	var cF []float64
	for i := 0; i < len(mcBands); i++ {
		if yAbs >= mcBands[i] {
			cF = mc2ll[i]
			break
		}
	}

	if cF == nil {
		return 0, 0, fmt.Errorf("coordinates out of range")
	}

	lng := cF[0] + cF[1]*xAbs
	latTemp := yAbs / cF[9]

	// 计算纬度多项式
	lat := cF[2] + cF[3]*latTemp +
		cF[4]*math.Pow(latTemp, 2) +
		cF[5]*math.Pow(latTemp, 3) +
		cF[6]*math.Pow(latTemp, 4) +
		cF[7]*math.Pow(latTemp, 5) +
		cF[8]*math.Pow(latTemp, 6)

	if x < 0 {
		lng *= -1
	}
	if y < 0 {
		lat *= -1
	}

	return lng, lat, nil
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
