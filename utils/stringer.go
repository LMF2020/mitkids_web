package utils

// 由于memcache key长度的限制，截取前十位与后十位保存在cache里，
// memcaheKey = jwt的shortkey
// memcacheValue = jwt的value
// memcacheTTL = 触发 /user/logout时，token剩余的expired time
// jwt token 过期会自动从memcached 移除
func ShortJwt(source string, first, last int) string {
	if len(source) < first+last {
		return ""
	}
	s := []rune(source)
	fstr := s[0:first]
	lstr := s[len(s)-last:]
	return string(append(fstr, lstr...))
}
