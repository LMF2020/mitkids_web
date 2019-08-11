package consts

const CodeExpiry = 180 // 验证码有效期3分钟
const CodeRegPrefix = "_reg_%s" // 注册用户, 验证码key前缀
const CodeLoginPrefix = "_lo_%s" // 登录用户, 验证码key前缀
const CodeForgetPassPrefix = "_fgp_%s" // 找回密码, 验证码key前缀

// 验证码类型
const CodeTypeReg = 1  // 验证码注册
const CodeTypeLogin = 2  // 验证码登录
const CodeTypeForgetPass = 3  // 验证码找回密码

const MaxBoundValueOfSearchRooms = 5 // km