package consts

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

	BOOK_MIN_UNIT = 1
	BOOK_MAX_UNIT = 8
)

var BOOK_LEVEL_SET = map[uint]uint{BookLevel1: 0, BookLevel2: 0, BookLevel3: 0}
