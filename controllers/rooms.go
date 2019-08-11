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

var wg sync.WaitGroup

// 根据geo查询可用的教室
func RoomsBoundsQueryHandler(c *gin.Context) {
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

	var rooms []model.Room
	if rooms, err = s.ListRoomByStatus(consts.RoomAvailable); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if rooms == nil {
		api.Success(c, []model.Room{}) // 没有数据
		return
	}

	// 如果数据有数据
	queue := make(chan model.Room, len(rooms))
	// 处理匹配的数据
	for _, room := range rooms {
		wg.Add(1)
		go MatchRoom(queue, room, strLat, strLng)
	}
	wg.Wait()
	close(queue)

	// 清空切片，重新添加
	rooms = []model.Room{}
	for r := range queue {
		rooms = append(rooms, r)
	}

	api.Success(c, rooms)
}

// 查询匹配5公里范围内的教室
func MatchRoom(queue chan model.Room, room model.Room, lan, lng float64) {
	defer wg.Done()
	distance := geo.GetDistance(room.Lat, lan, room.Lng, lng)
	if distance < consts.MaxBoundValueOfSearchRooms {
		queue <- room
	}
}
