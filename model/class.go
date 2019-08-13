package model

import (
	"mitkid_web/consts"
	"time"
)

// 班級
type Class struct {
	ClassId       string    `json:"class_id" form:"class_id" gorm:"primary_key"`      // 6位班级编号
	ClassName     string    `json:"class_name" form:"class_name" validate:"required"` // 6位班级编号
	ForeTeacherId string    `json:"fore_teacher_id" form:"fore_teacher_id" `          // 6位外教老師编号
	TeacherId     string    `json:"teacher_id" form:"teacher_id" `                    // 6位中教老師编号
	RoomId        string    `json:"room_id" form:"room_id" `                          // 上课教室 ID
	BookLevel     uint      `json:"book_level" form:"book_level" `                    // 课程级别
	BookPlan      string    `json:"book_plan" form:"book_plan"`                       // 课程系列(套餐)
	Status        uint      `json:"status" form:"status" `                            // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	ChildNumber   int       `json:"child_number" form:"child_number" `                // 当前报名人数
	Capacity      int       `json:"capacity" form:"capacity" `                        // 班級计划人数
	StartTime     time.Time `json:"start_time" form:"start_time"`                     // 开班时间
	EndTime       time.Time `json:"end_time" form:"end_time"`                         // 闭班时间
	CreatedAt     time.Time `json:"create_at" form:"create_at"`                       // 创建时间
	UpdatedAt     time.Time `json:"update_at" form:"update_at"`                       // 更新时间
	Childs        []string  `json:"childs" form:"childs" gorm:"-"`                    // 学生id列表
}

// 定义表名
func (class *Class) TableName() string {
	return consts.TABLE_CLASS
}
