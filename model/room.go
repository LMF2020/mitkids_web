package model

import (
	"mitkid_web/consts"
	"time"
)

// 教室
type Room struct {
	RoomId    int       `json:"room_id" gorm:"primary_key;auto_increment" sql:"type:int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY"` // 教室id自增
	RoomName  string    `json:"name" form:"name" gorm:"column:name" validate:"required"`
	Type      uint      `json:"type" form:"type" validate:"required,gte=1,lte=2"`           // 1机构提供 2合作家庭提供
	AccountId string    `json:"account_id" form:"account_id"`                               // 属于合作家庭需要绑定合作家庭编号
	Capacity  int       `json:"capacity" form:"capacity" validate:"required,gte=2,lte=500"` // 教室容纳学生数量
	Status    uint      `json:"status" form:"status" validate:"required,gte=1,lte=2"`       // 教室是否可用 1可用 2不可用
	ImageUrl  string    `json:"image_url" form:"image_url"`                                 //上传教室图片链接
	Lat       float64   `json:"lat" form:"lat" validate:"required"`                         // 经度, 以此坐标查询在此上课的所有班级
	Lng       float64   `json:"lng" form:"lng" validate:"required"`                         // 纬度, 以此坐标查询在此上课的所有班级
	CreatedAt time.Time `json:"create_at" `
	UpdatedAt time.Time `json:"update_at" `
	GeoAddr   string    `json:"geo_addr" form:"geo_addr" validate:"required"` // 地图认证的经纬度的地点名称
	Address   string    `json:"address" form:"address"`                       // 教师的具体地点(地图标记地址的补充)
	Phone     string    `json:"phone" form:"phone"`
}

// 定义表名
func (room *Room) TableName() string {
	return consts.TABLE_MK_ROOM
}

type RoomPageInfo struct {
	PageNumber int         `json:"page_number" form:"page_number" validate:"required" gorm:"-"`
	PageSize   int         `json:"page_size" form:"page_size" validate:"required" gorm:"-"`
	PageCount  int         `json:"page_count" gorm:"-"`
	TotalCount int         `json:"total_count" gorm:"-"`
	Results    interface{} `json:"results" gorm:"-"`
	Type       int         `form:"type" json:"-"`
	Status     int         `form:"status" json:"-"`
}

// 定义表名
func (r *RoomPageInfo) TableName() string {
	return consts.TABLE_MK_ROOM
}
