package model

import "time"

type class struct {
	classId       string    `json:"cliass_id" form:"cliass_id" gorm:"primary_key"`
	foreTeacherId string    `json:"fore_teacher_id" form:"fore_teacher_id"`
	teacherId     string    `json:"teacher_id" form:"teacher_id"`
	roomId        string    `json:"room_id" form:"room_id"`
	status        int       `json:"status" form:"status"`
	childNumber   int       `json:"child_number" form:"child_number"`
	bookCode      int       `json:"book_code" form:"book_code"`
	endTime       time.Time `json:"end_time" form:"end_time"`
	startTime     time.Time `json:"start_time" form:"start_time"`
	name          string    `json:"name" form:"name"`
	childIds      []string  `json:"child_ids" form:"child_ids"`
}

// 定义表名
func (class *class) TableName() string {
	return "mk_class"
}
