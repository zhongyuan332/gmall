package sessions

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

// Manager 是全局会话管理器
var Manager *sessions.Sessions

// Initialize 初始化会话管理器
func Initialize(cookieName string, expires time.Duration) {
	Manager = sessions.New(sessions.Config{
		Cookie:  cookieName,
		Expires: expires,
	})
}

// Get 从上下文获取会话
func Get(ctx iris.Context) *sessions.Session {
	return Manager.Start(ctx)
}

// Destroy 销毁当前会话
func Destroy(ctx iris.Context) {
	Manager.Destroy(ctx)
}
