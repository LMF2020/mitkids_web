package dao

import (
	"errors"
	"fmt"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"strings"
	"time"
)

// 查询学生最近要上的(N)节课
func (d *Dao) ListScheduledOccurringClass(classId, scheduledTimeOrder string, occurStatus, n int) (classOccurList []model.OccurClassPoJo, err error) {

	sql := `SELECT 
			  coo.class_id,
			  c.teacher_id,
			  c.fore_teacher_id,
			  c.book_level,
              c.class_name,
			  at_1.account_name AS teacher_name,
			  at_2.account_name AS fore_teacher_name,
			  rm.name AS room_name,
			  coo.book_code,
			  bk.book_name,
			  bk.book_link,
			  coo.occurrence_status AS status,
			  coo.occurrence_time
			FROM
			  mk_class_occurrence coo 
			  LEFT JOIN mk_class c 
				ON coo.class_id = c.class_id 
			  LEFT JOIN mk_room rm 
				ON rm.room_id = c.room_id 
			  LEFT JOIN mk_book bk 
				ON bk.book_code = coo.book_code 
			  LEFT JOIN mk_account at_1 
				ON at_1.account_id = c.teacher_id 
			  LEFT JOIN mk_account at_2 
				ON at_2.account_id = c.fore_teacher_id 
			WHERE c.class_id = ? 
			  AND coo.occurrence_status = ?
			ORDER BY coo.schedule_time %s
			LIMIT ?
			`
	sql = fmt.Sprintf(sql, scheduledTimeOrder)

	err = d.DB.Raw(sql, classId, occurStatus, n).Scan(&classOccurList).Error
	return
}

//
func (d *Dao) CountClassOccurs(classId string, occurStatus int) (count int, err error) {
	err = d.DB.Model(&model.ClassOccurrence{}).Where("class_id = ? and occurrence_status = ?", classId, occurStatus).Count(&count).Error
	return
}

func (d *Dao) CountClassOccursList(classIdArr []string, occurStatus int) (count int, err error) {
	err = d.DB.Model(&model.ClassOccurrence{}).Where("class_id in (?)  and occurrence_status = ?", classIdArr, occurStatus).Count(&count).Error
	return
}

// 分页查询上课历史
func (d *Dao) PageFinishedOccurrenceByClassIdArray(offset, pageSize int, classIdArr []string) (classOccurList []model.OccurClassPoJo, err error) {
	if classIdArr == nil || len(classIdArr) == 0 {
		classOccurList, err = nil, nil
		return
	}
	sqlStart := `SELECT 
			  coo.class_id,
			  coo.teacher_id,
			  coo.fore_teacher_id,
			  c.class_name,
			  c.book_level,
			  at_1.account_name AS teacher_name,
			  at_2.account_name AS fore_teacher_name,
			  rm.name AS room_name,
			  coo.book_code,
			  bk.book_name,
			  bk.book_link,
			  coo.occurrence_status AS status,
			  coo.schedule_time 
			FROM
			  mk_class_occurrence coo 
			  LEFT JOIN mk_class c 
				ON coo.class_id = c.class_id 
			  LEFT JOIN mk_room rm 
				ON rm.room_id = coo.room_id 
			  LEFT JOIN mk_book bk 
				ON bk.book_code = coo.book_code 
			  LEFT JOIN mk_account at_1 
				ON at_1.account_id = coo.teacher_id 
			  LEFT JOIN mk_account at_2 
				ON at_2.account_id = coo.fore_teacher_id 
			WHERE c.class_id in (`

	sqlEnd := `
				) 
				  AND coo.occurrence_status = ?
				ORDER BY coo.schedule_time DESC
				LIMIT ? OFFSET ?
				`
	sql := fmt.Sprintf("%s%s%s", sqlStart, strings.Join(classIdArr, ","), sqlEnd)
	err = d.DB.Raw(sql, consts.ClassOccurStatusFinished, pageSize, offset).Scan(&classOccurList).Error
	return
}

// 分页查询上课历史
func (d *Dao) PageFinishedOccurrenceByClassId(offset, pageSize int, classId string) (classOccurList []model.OccurClassPoJo, err error) {
	sql := `SELECT 
			  coo.class_id,
			  coo.teacher_id,
			  coo.fore_teacher_id,
  			  c.class_name,
			  c.book_level,
			  at_1.account_name AS teacher_name,
			  at_2.account_name AS fore_teacher_name,
			  rm.name AS room_name,
			  coo.book_code,
			  bk.book_name,
			  bk.book_link,
			  coo.occurrence_status AS status,
			  coo.schedule_time 
			FROM
			  mk_class_occurrence coo 
			  LEFT JOIN mk_class c 
				ON coo.class_id = c.class_id 
			  LEFT JOIN mk_room rm 
				ON rm.room_id = coo.room_id 
			  LEFT JOIN mk_book bk 
				ON bk.book_code = coo.book_code 
			  LEFT JOIN mk_account at_1 
				ON at_1.account_id = coo.teacher_id 
			  LEFT JOIN mk_account at_2 
				ON at_2.account_id = coo.fore_teacher_id 
			WHERE c.class_id = ? 
			  AND coo.occurrence_status = ?
			ORDER BY coo.schedule_time DESC
			LIMIT ? OFFSET ?
			`
	err = d.DB.Raw(sql, classId, consts.ClassOccurStatusFinished, pageSize, offset).Scan(&classOccurList).Error
	return
}

// 班级课程日历：包含课程是否结束的状态
func (d *Dao) ListOccurrenceCalendar(classId string) (classOccurList []model.OccurClassPoJo, err error) {
	sql := `SELECT 
			  coo.class_id,
			  coo.teacher_id,
			  coo.fore_teacher_id,
			  c.book_level,
			  at_1.account_name AS teacher_name,
			  at_2.account_name AS fore_teacher_name,
			  rm.name AS room_name,
			  coo.book_code,
			  bk.book_name,
			  bk.book_link,
			  coo.occurrence_status AS status,
			  coo.schedule_time,
              coo.occurrence_time 
			FROM
			  mk_class_occurrence coo 
			  LEFT JOIN mk_class c 
				ON coo.class_id = c.class_id 
			  LEFT JOIN mk_room rm 
				ON rm.room_id = coo.room_id 
			  LEFT JOIN mk_book bk 
				ON bk.book_code = coo.book_code 
			  LEFT JOIN mk_account at_1 
				ON at_1.account_id = coo.teacher_id 
			  LEFT JOIN mk_account at_2 
				ON at_2.account_id = coo.fore_teacher_id 
			WHERE c.class_id = ? 
			ORDER BY coo.schedule_time ASC
			`
	err = d.DB.Raw(sql, classId).Scan(&classOccurList).Error
	return
}

//添加课程
func (d *Dao) AddOccurrences(classId string, cOs *[]model.ClassOccurrence) (err error) {
	sql := genAddOccurrenceSql(classId, cOs)
	if err = d.DB.Exec(sql).Error; err != nil {
		log.Logger.Errorf("添加课程到课表失败：classId:%s，课程:%s,err:%s", classId, cOs, err)
		return errors.New("添加课程到课表失败")
	}
	return nil
}

func genAddOccurrenceSql(classId string, cOs *[]model.ClassOccurrence) (sql string) {
	sql = "INSERT INTO `mk_class_occurrence`(`class_id`, `occurrence_time`, `book_code`, `schedule_time`, `occurrence_status`, `create_at`, `updated_at`) VALUES "
	// 循环data数组,组合sql语句
	insertValuesFmt := "('%s','%s','%s', '%s', %d, NOW(), NOW()),"
	insertValuesFmt = fmt.Sprintf(insertValuesFmt, classId, "%s", "%s", "", consts.ClassOccurStatusNotStart)
	for _, cO := range *cOs {
		sql += fmt.Sprintf(insertValuesFmt, cO.OccurrenceTime, cO.BookCode)
	}
	sql = sql[0:len(sql)-1] + ";"
	return
}

const GetClassOccurrencesByClassId_sql = "select occurrence_time from mk_class_occurrence where class_id=?"

func (d *Dao) GetClassOccurrencesByClassId(classId string) (occurrences *[]time.Time, err error) {
	//rows, err :=d.DB.Table(consts.TABLE_CLASS_OCCURRENCE).Where("class_id = ?",classId).Select("occurrence_time").Rows()
	//defer rows.Close()
	//for rows.Next() {
	//	var o time.Time
	//	if err = rows.Scan(&o); err != nil {
	//		log.Logger.Error("row.Scan() error(%v)", err)
	//		return
	//	}
	//	occurrences = append(occurrences, o)
	//}
	occurrences = new([]time.Time)
	//d.DB.Raw(GetClassOccurrencesByClassId_sql, classId).Scan(occurrences)
	err = d.DB.Table(consts.TABLE_CLASS_OCCURRENCE).Where("class_id = ?", classId).Pluck("occurrence_time", occurrences).Error
	return
}

const EndClassOccurrClassOccurrencesByDateTimeSql = `UPDATE mk_class_occurrence co,
														mk_class c 
														SET co.occurrence_status = 2 
														WHERE
															co.class_id = c.class_id 
															and co.occurrence_status=1
														AND co.occurrence_time < ? 
														OR (
														co.occurrence_time = ? 
														AND c.end_time < ?)`

func (d *Dao) EndClassOccurrClassOccurrencesByDateTimeSql(datetime *time.Time) error {
	date := datetime.Format("2006-01-02 00:00:00")
	time := datetime.Format("15:04:05")
	return d.DB.Exec(EndClassOccurrClassOccurrencesByDateTimeSql, date, date, time).Error
}
