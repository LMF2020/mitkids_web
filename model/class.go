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

// 学生报名班级的关联表: 管理学生加入课堂状态
type JoinClass struct {
	ClassId         string    `json:"class_id" form:"class_id" gorm:"primary_key"`  // 6位班级编号
	AccountId     	string    `json:"account_id" form:"account_id" gorm:"primary_key"`
	Status        	uint      `json:"status" form:"status" `                            // 申请加入班级状态(1:申请中,2:申请加入成功,3:申请加入失败)
}

// 定义表名
func (class *JoinClass) TableName() string {
	return consts.TABLE_JOIN_CLASS
}

// 课程表：管理员或合作家庭在日历分配的课程表
type ClassOccurrence struct {
	ClassId         string    `json:"class_id" form:"class_id" gorm:"primary_key"` // 6位班级编号
	ScheduleTime         time.Time `json:"schedule_time" form:"schedule_time"`          // 日历规定的上课时间
	OccurrenceTime         time.Time `json:"occurrence_time" form:"occurrence_time"`          // 实际上课时间                                   // 闭班时间
	ForeTeacherId   string    `json:"fore_teacher_id" form:"fore_teacher_id" validate:"required"`          // 6位外教老師编号
	TeacherId       string    `json:"teacher_id" form:"teacher_id" validate:"required"`                    // 6位中教老師编号
	BookCode        string    `json:"book_code" form:"book_code"`             // 单节课的代码
	OccurrenceStatus          uint      `json:"occurrence_status" form:"occurrence_status" validate:"required"`    //  该课是否结束
	RoomId    string    `json:"room_id" form:"room_id"  validate:"required"` // 教室id
	ChildNumber     int       `json:"child_number" form:"child_number" validate:"required"`                // 上课实际人数
	Duration     int       `json:"duration" form:"duration"`   // 上课时长
	CreatedAt       time.Time `json:"create_at" form:"create_at"`                                          // 创建时间
	UpdatedAt       time.Time `json:"update_at" form:"update_at"`                                          // 更新时间
}

// 定义表名
func (class *ClassOccurrence) TableName() string {
	return consts.TABLE_CLASS_OCCURRENCE
}