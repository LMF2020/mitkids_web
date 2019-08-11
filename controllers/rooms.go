package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils/geo"
	"net/http"
	"strconv"
	"sync"
)

// 根据geo查询可用的教室
func RoomsBoundsQueryHandler(c *gin.Context) {
	var wg sync.WaitGroup
	lat := c.PostForm("lat")
	lng := c.PostForm("lng")

	var strLat, strLng float64
	var err error

	// 检查参数合法性
	if strLat, err = strconv.ParseFloat(lat, 64); err != nil {
		api.Fail(c, errorcode.INVALID_GEO, "参数无效")
		return
	}
	if strLng, err = strconv.ParseFloat(lng, 64); err != nil {
		api.Fail(c, errorcode.INVALID_GEO, "参数无效")
		return
	}

	var rooms [] model.Room
	if rooms, err = s.ListRoomByGeoAndStatus(strLat, strLng, consts.RoomAvailable); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if rooms == nil {
		api.Success(c, [] model.Room{})  // 没有数据
		return
	} else {

		var queue chan model.Room
		// 处理匹配的数据
		for _, room := range rooms {
			wg.Add(1)
			go MatchRoom(wg, queue, room, strLat, strLng)
		}
		wg.Wait()
		close(queue)
	}
}

// 查询匹配5公里内的教室
func MatchRoom(wg sync.WaitGroup, queue chan model.Room, room model.Room, lan, lng float64) {
	defer wg.Done()
	distance := geo.GetDistance(room.Lat, lan, room.Lng, lng)
	if distance < consts.MaxBoundValueOfSearchRooms {
		queue <- room
	}
}
