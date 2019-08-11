package main

import (
	"fmt"
	"mitkid_web/utils/geo"
	"reflect"
	"testing"
)

// 31.809407, 117.201134 // 习友路小学
// 31.831901, 117.139602 // 动漫基地
// 31.852291, 117.237059 // 萬象城
func TestGeo(t *testing.T) {
	lat1 := 31.809407
	lng1 := 117.201134
	lat2 := 31.831901
	lng2 := 117.139602
	fmt.Println(geo.GetDistance(lat1, lat2, lng1, lng2)) // 6km
	fmt.Println(reflect.TypeOf(lat1))                    // float64
}
