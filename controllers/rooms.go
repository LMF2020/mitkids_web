package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/geo"
	"net/http"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

// 查询匹配5公里范围内的教室
func MatchRoom(queue chan model.Room, room model.Room, lan, lng float64) {
	defer wg.Done()
	distance := geo.GetDistance(room.Lat, lan, room.Lng, lng)
	if distance < consts.MaxBoundValueOfSearchRooms {
		queue <- room
	}
}

func CreateRoom(c *gin.Context) {
	var from model.Room
	var err error
	if err = c.ShouldBind(&from); err == nil {
		if err = utils.ValidateParam(from); err == nil {
			exist, err := s.GetRoomByAddAndGeo(from.GeoAddr, from.Address)
			if err != nil {
				api.Failf(c, http.StatusBadRequest, "查询错误", err)
				return
			}
			if exist != nil {
				api.Failf(c, http.StatusBadRequest, "教室地址重复")
				return
			}
			if err = s.CreateRoom(&from); err == nil {
				api.Success(c, "新建教室成功")
				return
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
}

//获取教室
func GetRoomById(c *gin.Context) {
	idStr := c.PostForm("room_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Failf(c, http.StatusBadRequest, "参数错误 id:%s", idStr)
		return
	}
	if room, err := s.GetRoomById(id); err == nil {
		if room == nil {
			api.Failf(c, http.StatusBadRequest, "不存在这个教室id:%s", idStr)
			return
		}

		api.Success(c, room)
		return
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

//获取教室
func DeleteRoomById(c *gin.Context) {
	idStr := c.PostForm("room_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Failf(c, http.StatusBadRequest, "参数错误 id:%s", idStr)
		return
	}
	if room, err := s.GetRoomById(id); err == nil {
		if room == nil {
			api.Failf(c, http.StatusBadRequest, "不存在这个教室id:%s", idStr)
			return
		}
		if err = s.DeleteRoomById(id); err == nil {
			api.Success(c, "删除教室成功")
			return
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

//获取教室
func UpdateRoomById(c *gin.Context) {
	idStr := c.PostForm("room_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Failf(c, http.StatusBadRequest, "参数错误 id:%s", idStr)
		return
	}
	if room, err := s.GetRoomById(id); err == nil {
		if room == nil {
			api.Failf(c, http.StatusBadRequest, "不存在这个教室id:%s", idStr)
			return
		}
		if err = c.ShouldBind(room); err == nil {
			if err = utils.ValidateParam(room); err == nil {
				if err = s.UpdateRoom(room); err == nil {
					api.Success(c, "更新教室成功")
					return
				}
			}
		}
		api.Success(c, room)
		return
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}

func ListRoomWithQueryByPage(c *gin.Context) {
	var pageInfo model.RoomPageInfo
	var err error
	//pageInfo.Type=consts.RoomTyep_1
	//pageInfo.Status=consts.RoomAvailable
	if err = c.ShouldBind(&pageInfo); err == nil {
		if err = utils.ValidateParam(pageInfo); err == nil {
			pn, ps := pageInfo.PageNumber, pageInfo.PageSize
			if pn < 0 {
				pn = 1
			}
			if ps <= 0 {
				ps = consts.DEFAULT_PAGE_SIZE
			}
			query := c.PostForm("query")
			totalRecords, err := s.CountRoomWithQuery(&pageInfo, query)
			if err != nil {
				api.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if totalRecords == 0 {
				api.Success(c, pageInfo)
				return
			}
			pageCount := totalRecords / ps
			if totalRecords%ps > 0 {
				pageCount++
			}
			if pn > pageCount {
				pn = pageCount
			}
			pageInfo.PageCount = pageCount
			pageInfo.TotalCount = totalRecords
			if accounts, err := s.ListRoomWithQueryByPage(pn, ps, &pageInfo, query); err == nil {
				pageInfo.Results = accounts
				api.Success(c, pageInfo)
				return
			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return

}
