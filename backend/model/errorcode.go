package model

type errorCode struct {
	SUCCESS      int
	ERROR        int
	NotFound     int
	LoginError   int
	LoginExist   int
	LoginExpired int
}

// ErrorCode 错误码
var ErrorCode = errorCode{
	SUCCESS:      0,
	ERROR:        1,
	NotFound:     404,
	LoginError:   1000, //用户名或密码错误
	LoginExist:   1001, //用户已存在
	LoginExpired: 1002, //登录过期
}
