package model

import "time"

// 班級
type Class struct {
	ClassId         string    `json:"class_id" form:"class_id" gorm:"primary_key"`                         // 6位班级编号
	ForeTeacherId   string    `json:"fore_teacher_id" form:"fore_teacher_id" validate:"required"`          // 6位外教老師编号
	TeacherId       string    `json:"teacher_id" form:"teacher_id" validate:"required"`                    // 6位中教老師编号
	RoomId          string    `json:"room_id" form:"room_id" validate:"required"`                          // 上课教室 ID
	BookLevel       uint      `json:"book_level" form:"book_level" validate:"required"`                    // 课程级别
	BookPlan        string    `json:"book_plan" form:"book_plan"`                                          // 课程系列(套餐)
	Status          uint      `json:"status" form:"status" validate:"required"`                            // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	ChildNumber     int       `json:"child_number" form:"child_number" validate:"required"`                // 当前报名人数
	Capacity        int       `json:"capacity" form:"capacity" validate:"required"`                        // 班級计划人数
	StartTime       time.Time `json:"start_time" form:"start_time"`                                        // 开班时间
	EndTime         time.Time `json:"end_time" form:"end_time"`                                            // 闭班时间
	CreatedAt       time.Time `json:"create_at" form:"create_at"`                                          // 创建时间
	UpdatedAt       time.Time `json:"update_at" form:"update_at"`                                          // 更新时间
	TeacherName     string    `json:"teacher_name" form:"teacher_name" gorm:"teacher_name"`                // 中教姓名
	ForeTeacherName string    `json:"fore_teacher_name" form:"fore_teacher_name" gorm:"fore_teacher_name"` // 外教姓名
}

// 定义表名
func (class *Class) TableName() string {
	return "mk_class"
}
