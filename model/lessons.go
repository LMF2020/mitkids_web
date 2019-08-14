package model

import (
	"mitkid_web/consts"
	"time"
)

// 针对于每一节课
// 学生每节课的课堂表: 每次上完课，教师端更新上课状态，并给学生评，该表也可以统计学生是否缺席的情况
type Lessons struct {
	BookCode  string    `json:"book_code" form:"book_code" gorm:"primary_key"`   // 单节课的代码
	StudentId string    `json:"student_id" form:"student_id" gorm:"primary_key"` // 8位学生编号
	TeacherId string    `json:"teacher_id" form:"teacher_id"`                    // 8位学生编号
	ClassId   string    `json:"class_id" form:"class_id" gorm:"primary_key"`     // 6位班级编号
	Status    uint      `json:"status" form:"status" `                           // 完成状态; 1未上课 2已上课
	Score     uint      `json:"score" form:"score" `                             // 学生给老师打分;由低到高:1-5分
	Comment   string    `json:"comment" form:"comment"`                          // 老师给学生评语
	CreatedAt time.Time `json:"create_at" form:"create_at"`                      // 创建时间
	UpdatedAt time.Time `json:"update_at" form:"update_at"`                      // 更新时间
}

// 定义表名
func (class *Lessons) TableName() string {
	return consts.TABLE_LESSONS
}
