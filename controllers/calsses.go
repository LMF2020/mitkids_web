package controllers

import (
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"net/http"
)

// 返回报文：
/**
{
  LV1: [
	{
	  class_id
	  class_name
	  room_name
      teacher
	}
  ],
  LV2: [
	{
	  class_id
	  class_name
	  room_name
	  teacher
	}
  ],
  LV3: [
	{
	  class_id
	  class_name
	  room_name
	  teacher
	}
  ]
}
*/
// 根据教室查询所有班级
func ClassesQueryByRoomIdHandler(c *gin.Context) {

	roomId := c.Param("roomId")

	if classes, err := s.ListAvailableClassesByRoomId(roomId); err != nil {
		api.Fail(c, http.StatusInternalServerError, "请求内部错误")
		return
	} else if classes == nil {
		api.Success(c, make(map[string]model.Class)) // 没有数据
		return
	} else {
		// 报文解析
		retJson := make(map[string][]model.Class)
		var LV1, LV2, LV3 []model.Class
		for _, item := range classes {
			switch item.BookLevel {
			case consts.BookLevel1:
				LV1 = append(LV1, item)
			case consts.BookLevel2:
				LV2 = append(LV2, item)
			case consts.BookLevel3:
				LV3 = append(LV3, item)
			}
		}

		if len(LV1) > 0 {
			retJson["LV1"] = LV1
		}

		if len(LV2) > 0 {
			retJson["LV2"] = LV2
		}

		if len(LV3) > 0 {
			retJson["LV3"] = LV3
		}

		api.Success(c, retJson)
	}

}

func CreateClass(c *gin.Context) {
	var formClass model.Class
	var err error
	if err = c.ShouldBind(&formClass); err == nil {
		if err = utils.ValidateParam(formClass); err == nil {
			if _, ok := consts.BOOK_LEVEL_SET[formClass.BookLevel]; !ok {
				api.Fail(c, http.StatusBadRequest, "无效的课程")
				return
			}
			if formClass.BookFromUnit > formClass.BookToUnit {
				api.Fail(c, http.StatusBadRequest, "课程开始单元不能大于结束单元")
				return
			}
			if formClass.BookFromUnit < consts.BOOK_MIN_UNIT || formClass.BookFromUnit > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程开始单元无效")
				return
			}
			if formClass.BookToUnit < consts.BOOK_MIN_UNIT || formClass.BookToUnit > consts.BOOK_MAX_UNIT {
				api.Fail(c, http.StatusBadRequest, "课程结束单元无效")
				return
			}
			formClass.ChildNumber = uint(len(formClass.Childs))
			if err = s.CreateClass(&formClass); err == nil {
				if formClass.ChildNumber != 0 {
					if err = s.AddChildsToClass(formClass.ClassId, formClass.Childs); err == nil {
						api.Success(c, "创建班级成功")
					}
				}

			}
		}
	}
	api.Fail(c, http.StatusBadRequest, err.Error())
	return
}
