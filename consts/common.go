package consts

import "fmt"

const (
	CodeExpiry           = 180       // 验证码有效期3分钟
	CodeRegPrefix        = "_reg_%s" // 注册用户, 验证码key前缀
	CodeLoginPrefix      = "_lo_%s"  // 登录用户, 验证码key前缀
	CodeForgetPassPrefix = "_fgp_%s" // 找回密码, 验证码key前缀

	// 验证码类型
	CodeTypeReg        = 1 // 验证码注册
	CodeTypeLogin      = 2 // 验证码登录
	CodeTypeForgetPass = 3 // 验证码找回密码

	MaxBoundValueOfSearchRooms = 6 // km

	BookLevel1 = 1
	BookLevel2 = 2
	BookLevel3 = 3

	DEFAULT_PAGE_SIZE = 10

	BOOK_MIN_UNIT         = 1
	BOOK_MAX_UNIT         = 3
	BOOK_UNIT_CLASS_COUNT = 8
	BOOK_PLAN_FMT         = "%s%d_%d单元"

	URL_CHILD_LOGOUT   = "/api/child/logout"
	URL_ADMIN_LOGOUT   = "/api/admin/logout"
	URL_TEACHER_LOGOUT = "/api/teacher/logout"

	JWT_SECRETS = "MITSECRET2019"
	JWT_VENDOR  = "MITVENDOR2019"

	STATUS_CHILD_CLASS_MISSED   = 1
	STATUS_CHILD_CLASS_ATTENDED = 2
)

var (
	BOOK_LEVEL_SET      = map[uint]string{BookLevel1: "初级", BookLevel2: "中级", BookLevel3: "高级"}
	URL_LOGOUT_API_LIST = fmt.Sprintf("%s_%s_%s", URL_CHILD_LOGOUT, URL_ADMIN_LOGOUT, URL_TEACHER_LOGOUT)
	REGEX_TEACHER_API   = "^/api/teacher/.*"
	REGEX_CHILD_API     = "^/api/child/.*"
	REGEX_CORP_API      = "^/api/corp/.*"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
)
