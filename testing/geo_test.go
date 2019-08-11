package main

import (
	"fmt"
	"mitkid_web/utils/geo"
	"reflect"
	"testing"
)
// 31.810512, 117.202707 // 习友路小学
// 31.831901, 117.139602 // 动漫基地
func TestGeo(t *testing.T) {
	lat1 := 31.810512
	lng1 := 117.202707
	lat2 := 31.831901
	lng2 := 117.139602
	fmt.Println(geo.GetDistance(lat1, lat2, lng1, lng2)) // 6km
	fmt.Println(reflect.TypeOf(lat1)) // float64
}