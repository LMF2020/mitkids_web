package errorcode

const (
	// user :1000-1100
	USER_ALREADY_EXIS = 1000
	USER_NOT_EXIS     = 1001

	// common ：1101-1200
	VERIFY_CODE_ERR = 1101
	INVALID_GEO     = 1102

	// token 过期，不能刷新获取
	ErrExpiredToken = 401001
	// header token为空
	ErrEmptyAuthHeader = 401002
	// token无效
	ErrInvalidAuthHeader = 401003
	// cookie token为空
	ErrEmptyCookieToken = 401004
	// token签名算法错误
	ErrInvalidSigningAlgorithm = 401005
	// 其他原因的错误
	ErrOtherCase = 401000
	// class ：1201-1300

)
