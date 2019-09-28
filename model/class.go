package model

import (
	"mitkid_web/consts"
	"mitkid_web/utils"
	"time"
)

// 班級
type Class struct {
	ClassId       string        `json:"class_id" form:"class_id" gorm:"primary_key"`      // 6位班级编号
	ClassName     string        `json:"class_name" form:"class_name" validate:"required"` // 6位班级名称
	ForeTeacherId string        `json:"fore_teacher_id" form:"fore_teacher_id" `          // 6位外教老師编号
	TeacherId     string        `json:"teacher_id" form:"teacher_id" `                    // 6位中教老師编号
	RoomId        string        `json:"room_id" form:"room_id" validate:"required"`       // 上课教室 ID
	BookLevel     uint          `json:"book_level" form:"book_level" `                    // 课程级别
	Status        uint          `json:"status" form:"status" `                            // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	ChildNumber   uint          `json:"child_number" form:"child_number" `                // 当前报名人数
	Capacity      uint          `json:"capacity" form:"capacity" validate:"required"`     // 班級计划人数
	StartTime     utils.RawTime `json:"start_time" form:"start_time" validate:"required"` // 课程开始时间
	EndTime       utils.RawTime `json:"end_time" form:"end_time" validate:"required"`     // 课程结束时间
	Childs        []string      `json:"childs" form:"childs" gorm:"-"`                    // 学生id列表
	ChildNames    []string      `json:"child_names" form:"_" gorm:"-"`
	BookFromUnit  uint          `json:"book_from_unit" form:"book_from_unit"  validate:"required" ` // 课程 开始单元
	BookToUnit    uint          `json:"book_to_unit" form:"book_to_unit"  validate:"required" `     // 课程 结束单元
	Occurrences   []time.Time   `json:"occurrences" form:"occurrences" gorm:"-" validate:"required" time_format:"2006-01-02"`
	Weeks         string        `json:"weeks" form:"weeks" `                                   // 每周具体上课天数
	StartDate     time.Time     `json:"start_date" form:"start_date" time_format:"2006-01-02"` // 课程开始日期
	CreatedAt     time.Time     `json:"create_at" form:"create_at"`                            // 创建时间
	UpdatedAt     time.Time     `json:"update_at" form:"update_at"`                            // 更新时间

	// 从roomid关联过来的json字段
	GeoAddr   string    `json:"geo_addr" form:"geo_addr" gorm:"-"` 								// 地图认证的经纬度的地点名称
	Address   string    `json:"address" form:"address" gorm:"-"`                        		// 教师的具体地点(地图标记地址的补充)
	TeacherName string	`json:"teacher_name" form:"teacher_name" gorm:"-"`     					// 中教姓名
	ForeTeacherName string	`json:"fore_teacher_name" form:"fore_teacher_name" gorm:"-"`     	// 外教姓名

}

// 定义表名
func (class *Class) TableName() string {
	return consts.TABLE_CLASS
}

// 学生报名班级的关联表: 管理学生加入课堂状态
type JoinClass struct {
	ClassId   string `json:"class_id" form:"class_id" gorm:"primary_key" validate:"required"` // 6位班级编号
	AccountId string `json:"account_id" form:"account_id" gorm:"primary_key" validate:"required"`
	Status    int    `json:"status" form:"status"`
}

// 定义表名
func (class *JoinClass) TableName() string {
	return consts.TABLE_JOIN_CLASS
}

// 课程表：管理员或合作家庭在日历分配的课程表
type ClassOccurrence struct {
	ClassId          string    `json:"class_id" form:"class_id" gorm:"primary_key"`                    // 6位班级编号
	ScheduleTime     time.Time `json:"schedule_time" form:"schedule_time"`                             // 日历规定的上课时间
	OccurrenceTime   time.Time `json:"occurrence_time" form:"occurrence_time"`                         // 实际上课时间                                   // 闭班时间
	ForeTeacherId    string    `json:"fore_teacher_id" form:"fore_teacher_id" validate:"required"`     // 6位外教老師编号
	TeacherId        string    `json:"teacher_id" form:"teacher_id" validate:"required"`               // 6位中教老師编号
	BookCode         string    `json:"book_code" form:"book_code"`                                     // 单节课的代码
	OccurrenceStatus uint      `json:"occurrence_status" form:"occurrence_status" validate:"required"` //  该课是否结束
	RoomId           string    `json:"room_id" form:"room_id"  validate:"required"`                    // 教室id
	ChildNumber      int       `json:"child_number" form:"child_number" validate:"required"`           // 上课实际人数
	Duration         int       `json:"duration" form:"duration"`                                       // 上课时长
	CreatedAt        time.Time `json:"create_at" form:"create_at"`                                     // 创建时间
	UpdatedAt        time.Time `json:"update_at" form:"update_at"`                                     // 更新时间
}

// 定义表名
func (class *ClassOccurrence) TableName() string {
	return consts.TABLE_CLASS_OCCURRENCE
}

type ChildClass struct {
	ClassId       string        `json:"class_id" form:"class_id" gorm:"primary_key"`      // 6位班级编号
	ClassName     string        `json:"class_name" form:"class_name" validate:"required"` // 6位班级编号
	ForeTeacherId string        `json:"fore_teacher_id" form:"fore_teacher_id" `          // 6位外教老師编号
	TeacherId     string        `json:"teacher_id" form:"teacher_id" `                    // 6位中教老師编号
	RoomId        string        `json:"room_id" form:"room_id" `                          // 上课教室 ID
	BookLevel     uint          `json:"book_level" form:"book_level" `                    // 课程级别
	BookPlan      string        `json:"book_plan" form:"book_plan"`                       // 课程系列(套餐)
	Status        uint          `json:"status" form:"status" `                            // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	ChildNumber   uint          `json:"child_number" form:"child_number" `                // 当前报名人数
	Capacity      uint          `json:"capacity" form:"capacity" validate:"required"`     // 班級计划人数
	StartTime     utils.RawTime `json:"start_time" form:"start_time" validate:"required"` // 课程开始时间
	EndTime       utils.RawTime `json:"end_time" form:"end_time" validate:"required"`     // 课程结束时间    // 闭班时间
	CreatedAt     time.Time     `json:"create_at" form:"create_at"`                       // 创建时间
	UpdatedAt     time.Time     `json:"update_at" form:"update_at"`                       // 更新时间
	StudentId     string        `json:"student_id" gorm:"student_id"`
}

// 班級
type ClassListItem struct {
	ClassId         string        `json:"class_id" form:"class_id" gorm:"primary_key"`                // 6位班级编号
	ClassName       string        `json:"class_name" form:"class_name" validate:"required"`           // 6位班级名称
	ForeTeacherId   string        `json:"fore_teacher_id" form:"fore_teacher_id" `                    // 6位外教老師编号
	TeacherId       string        `json:"teacher_id" form:"teacher_id" `                              // 6位中教老師编号
	ForeTeacherName string        `json:"fore_teacher_name" form:"fore_teacher_name" `                // 6位外教老師编号
	TeacherName     string        `json:"teacher_name" form:"teacher_name" `                          // 6位中教老師编号
	RoomId          string        `json:"room_id" form:"room_id" validate:"required"`                 // 上课教室 ID
	GeoAddr         string        `json:"geo_addr" form:"geo_addr"`                                   // 地图认证的经纬度的地点名称
	Address         string        `json:"address" form:"address"`                                     // 教师的具体地点(地图标记地址的补充)
	BookLevel       uint          `json:"book_level" form:"book_level" `                              // 课程级别
	Status          uint          `json:"status" form:"status" `                                      // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	ChildNumber     uint          `json:"child_number" form:"child_number" `                          // 当前报名人数
	Capacity        uint          `json:"capacity" form:"capacity" validate:"required"`               // 班級计划人数
	StartTime       utils.RawTime `json:"start_time" form:"start_time" validate:"required"`           // 课程开始时间
	EndTime         utils.RawTime `json:"end_time" form:"end_time" validate:"required"`               // 课程结束时间
	Childs          []string      `json:"childs" form:"childs" gorm:"-"`                              // 学生id列表
	BookFromUnit    uint          `json:"book_from_unit" form:"book_from_unit"  validate:"required" ` // 课程 开始单元
	BookToUnit      uint          `json:"book_to_unit" form:"book_to_unit"  validate:"required" `     // 课程 结束单元
	Occurrences     []time.Time   `json:"occurrences" form:"occurrences" gorm:"-" validate:"required" time_format:"2006-01-02"`
	Weeks           string        `json:"weeks" form:"weeks" `                                   // 每周具体上课天数
	StartDate       time.Time     `json:"start_date" form:"start_date" time_format:"2006-01-02"` // 课程开始日期
	CreatedAt       time.Time     `json:"create_at" form:"create_at"`                            // 创建时间
	UpdatedAt       time.Time     `json:"update_at" form:"update_at"`                            // 更新时间
}

type JoinClassItem struct {
	ClassId   string        `json:"class_id" form:"class_id" gorm:"primary_key"`           // 6位班级编号
	ClassName string        `json:"class_name" form:"class_name" validate:"required"`      // 6位班级名称
	Status    uint          `json:"status" form:"status" `                                 // 班級是否关闭(1:未开始,2:进行中,3:已结束)
	StartTime utils.RawTime `json:"start_time" form:"start_time" validate:"required"`      // 课程开始时间
	EndTime   utils.RawTime `json:"end_time" form:"end_time" validate:"required"`          // 课程结束时间
	Weeks     string        `json:"weeks" form:"weeks" `                                   // 每周具体上课天数
	StartDate time.Time     `json:"start_date" form:"start_date" time_format:"2006-01-02"` // 课程开始日期
	StudentId string        `json:"student_id" gorm:"student_id"`                          // 学生
}
