package model

import "time"

// 教室
type Room struct {
	RoomId    string    `json:"room_id" form:"room_id" gorm:"primary_key;auto_increment"` // 教室id自增
	RoomName  string    `json:"name" form:"name" gorm:"column:name"`
	Type      uint      `json:"type" form:"type" validate:"required"`                       // 1机构提供 2合作家庭提供
	AccountId string    `json:"account_id" form:"account_id"`                               // 属于合作家庭需要绑定合作家庭编号
	Capacity  int       `json:"capacity" form:"capacity" validate:"required,gte=2,lte=500"` // 教室容纳学生数量
	Status    uint      `json:"status" form:"status" validate:"required"`                   // 教室是否可用 1可用 2不可用
	ImageUrl  string    `json:"image_url" form:"image_url"`                                 //上传教室图片链接
	Lat       float64   `json:"lat" form:"lat" validate:"required"`                         // 经度, 以此坐标查询在此上课的所有班级
	Lng       float64   `json:"lng" form:"lng" validate:"required"`                         // 纬度, 以此坐标查询在此上课的所有班级
	CreatedAt time.Time `json:"create_at" form:"create_at"`
	UpdatedAt time.Time `json:"update_at" form:"update_at"`
	GeoAddr   string    `json:"geo_addr" form:"geo_addr" validate:"required"` // 地图认证的经纬度的地点名称
	Address   string    `json:"address" form:"address"`                       // 教师的具体地点(地图标记地址的补充)
}

// 定义表名
func (room *Room) TableName() string {
	return "mk_room"
}
