package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/logger"
	model "github.com/zhongyuan332/gmall/backend/model"
	"github.com/zhongyuan332/gmall/sessions"
	"time"
)

// Session keys
const (
	UserIDKey     = "UserID"
	UsernameKey   = "Username"
	IsLoggedInKey = "IsLoggedIn"
	LastAccessKey = "LastAccess"
)

// Authentication 授权
func Authentication(ctx iris.Context) {
	logger.Debugf("admin authentication start...")

	session := sessions.Get(ctx)
	isLoggedIn, _ := session.GetBoolean(IsLoggedInKey)

	// 检查是否已登录
	if !isLoggedIn {
		ctx.Redirect("/login", iris.StatusFound)
		return
	}

	// 检查会话是否过期（例如，如果用户在30分钟内没有活动）
	lastAccess := session.GetInt64Default(LastAccessKey, 0)
	now := time.Now().Unix()
	inactiveFor := now - lastAccess

	// 如果不活动超过30分钟，则要求重新登录
	const maxInactiveTime int64 = 30 * 60 // 30分钟
	if inactiveFor > maxInactiveTime {
		// 清除会话并重定向到登录页面
		session.Delete(IsLoggedInKey)
		ctx.Redirect("/login?expired=true", iris.StatusFound)
		return
	}

	// 更新最后访问时间
	session.Set(LastAccessKey, now)

	// 将用户信息添加到context中
	userID := session.GetInt64Default(UserIDKey, 0)
	username := session.GetStringDefault(UsernameKey, "")

	ctx.Values().Set(UserIDKey, userID)
	ctx.Values().Set(UsernameKey, username)
	ctx.Values().Set(IsLoggedInKey, true)

	ctx.Next()
}

// ShowLogin 显示登录页面
func ShowLogin(ctx iris.Context) {
	// 检查是否因会话过期而重定向
	expired := ctx.URLParamExists("expired")

	// 传递给模板的数据
	data := iris.Map{
		"Title":   "登录系统",
		"Expired": expired,
	}

	ctx.View("auth/login.html", data)
}

// Login 处理登录请求
func Login(ctx iris.Context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	userService := &UserService{model.DB.DB()}
	// 简单的用户验证逻辑，实际应用中应该从数据库验证
	user, err := userService.VerifyPassword(username, password)
	if err != nil {
		// 登录失败
		ctx.ViewData("error", "用户名或密码错误")
		ctx.ViewData("username", username)
		ctx.View("auth/login.html")
		return
	}

	// 登录成功，设置会话
	session := sessions.Get(ctx)
	session.Set(UserIDKey, user.ID)
	session.Set(UsernameKey, user.Username)
	session.Set(IsLoggedInKey, true)
	session.Set(LastAccessKey, time.Now().Unix())

	// 重定向到首页或上一个请求的页面
	returnURL := ctx.FormValue("returnUrl")
	if returnURL == "" {
		returnURL = "/"
	}

	ctx.Redirect(returnURL, iris.StatusFound)
}

// Logout 处理登出请求
func Logout(ctx iris.Context) {
	session := sessions.Get(ctx)

	// 清除会话
	session.Delete(UserIDKey)
	session.Delete(UsernameKey)
	session.Delete(IsLoggedInKey)
	session.Delete(LastAccessKey)

	ctx.Redirect("/login", iris.StatusFound)
}
