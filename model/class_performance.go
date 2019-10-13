package model

import (
	"mitkid_web/consts"
	"mitkid_web/utils"
	"time"
)

// 学生课堂评分
type ClassPerformance struct {
	AccountId string    `json:"account_id" form:"account_id" gorm:"primary_key"` // 8位学生编号
	ClassId   string    `json:"class_id" form:"class_id" gorm:"primary_key"`     // 6位班级编号
	ClassDate string `json:"class_date" form:"class_date" gorm:"primary_key"` 	 // 上课日期
	TeacherId string    `json:"teacher_id" form:"teacher_id"`                    // 6位中教编号
	Status    int       `json:"status" form:"status" `                           // 完成状态; 1未上课 2已上课
	Comment   string    `json:"comment" form:"comment"`                          // 老师给学生评语
	Option1   string    `json:"option1" form:"option1"`                          // 评分大项1：小项星星的数量以逗号分隔 (2,3,2,2,5)
	Option2   string    `json:"option2" form:"option2"`                          // 评分大项2
	Option3   string    `json:"option3" form:"option3"`                          // 评分大项3
	Option4   string    `json:"option4" form:"option4"`                          // 评分大项4
	//Option5   string    `json:"option5" form:"option5"`                          // 评分大项5
	CreatedAt time.Time `json:"created_at" form:"created_at"`                      // 创建时间
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`                      // 更新时间
}

// 定义表名
func (class *ClassPerformance) TableName() string {
	return consts.TABLE_CLASS_PERFORMANCE
}

// 课表记录
type ClassRecordItem struct {
	ClassId         string        `json:"class_id" form:"class_id" `                  // 6位班级编号
	ClassName       string        `json:"class_name" form:"class_name"`               // 所在班级名称
	TeacherId       string        `json:"teacher_id" form:"teacher_id"`               // 中教编号
	ForeTeacherId   string        `json:"fore_teacher_id" form:"fore_teacher_id"`     // 外教编号
	TeacherName     string        `json:"teacher_name" form:"teacher_name"`           // 中教姓名
	ForeTeacherName string        `json:"fore_teacher_name" form:"fore_teacher_name"` // 外教姓名
	BookLevel       uint          `json:"book_level" form:"book_level" `              // 课程级别
	BookCode        string        `json:"book_code" form:"book_code"`                 // 课本的代码
	BookName        string        `json:"book_name" form:"book_name"`                 // 课本的名称
	RoomName        string        `json:"room_name" form:"room_name"`                 // 上课教室
	ScheduleTime    utils.RawTime `json:"schedule_time" form:"schedule_time"`         // 计划上课时间
	OccurrenceTime  string     	   `json:"occurrence_time" form:"occurrence_time"`     // 实际上课时间 YYYY-MM-DD
	GeoAddr         string        `json:"geo_addr" form:"geo_addr"`                   // 地图认证的经纬度的地点名称
	Address         string        `json:"address" form:"address"`
	BookLink        string        `json:"book_link" form:"book_link"` // 课本预习链接
	Status          uint          `json:"status" form:"status" `      // 完成状态; 1未上课 2已上课
}
