package consts

const (
	AccountRoleTeacher         = 1 // 中教
	AccountRoleCorp            = 2 // 合作家庭
	AccountRoleChild           = 3 // 学生
	AccountRoleAdmin           = 4 //  管理员
	AccountRoleForeignTeacher  = 5 // 外教
	AccountRoleCorpWithTeacher = 6 // 合作家庭中教

	AccountStatusNormal uint = 2
	AccountStatusFail   uint = 1
	AccountStatusWait        = 0

	AccountTypePaid = 2
	AccountTypeFree = 1

	AccountLoginTypePass = 1
	AccountLoginTypeCode = 2

	AccountGenderMale = 1 // 男
	AccountGenderFemale = 2  // 女
)
