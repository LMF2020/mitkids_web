package dao

import (
	"errors"
	"fmt"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
	"time"
)

// 查询学生最近要上的(N)节课
func (d *Dao) ListClassOccurrence(classId, scheduledTimeOrder string, occurStatus, limit int) (classOccurList []model.OccurClassPoJo, err error) {

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
			ORDER BY coo.schedule_time %s
			LIMIT ?
			`
	sql = fmt.Sprintf(sql, scheduledTimeOrder)

	err = d.DB.Raw(sql, classId, occurStatus, limit).Scan(&classOccurList).Error
	return
}

//
func (d *Dao) CountOccurrence(classId string, occurStatus int) (count int, err error) {
	err = d.DB.Model(&model.ClassOccurrence{}).Where("class_id = ? and occurrence_status = ?", classId, occurStatus).Count(&count).Error
	return
}

// 分页查询上课历史
func (d *Dao) ListOccurrenceHisByPage(offset, pageSize int, classId string) (classOccurList []model.OccurClassPoJo, err error) {
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

func (d *Dao) GetClassOccurrencesByClassId(classId string) (occurrences *[]time.Time) {
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
	d.DB.Raw(GetClassOccurrencesByClassId_sql, classId).Scan(&occurrences)
	//d.DB.Table(consts.TABLE_CLASS_OCCURRENCE).Where("class_id = ?", classId).Pluck("occurrence_time", &occurrences)
	return
}
