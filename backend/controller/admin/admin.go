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
		logger.Warn("User not logged in, redirecting to login page")
		ctx.JSON(iris.Map{
			"errNo": model.ErrorCode.LoginExpired,
			"msg":   "登录已过期",
		})
		return
	}

	// 检查会话是否过期（例如，如果用户在30分钟内没有活动）
	lastAccess := session.GetInt64Default(LastAccessKey, 0)
	now := time.Now().Unix()
	inactiveFor := now - lastAccess

	// 如果不活动超过30分钟，则要求重新登录
	const maxInactiveTime int64 = 30 * 60 // 30分钟
	if inactiveFor > maxInactiveTime {
		logger.Warnf("Session expired for user %s, inactive for %d seconds", session.GetStringDefault(UsernameKey, ""), inactiveFor)
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

// Login 处理登录请求
func Login(ctx iris.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := ctx.ReadJSON(&loginRequest)
	if err != nil {
		// 处理JSON解析错误
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"errNo":   400,
			"message": "无效的JSON格式: " + err.Error(),
		})
		return
	}
	logger.Infof("Login attempt with username: %s, password:%v", loginRequest.Username, loginRequest.Password)
	userService := &UserService{model.DB.DB()}
	// 简单的用户验证逻辑，实际应用中应该从数据库验证
	user, err := userService.VerifyPassword(loginRequest.Username, loginRequest.Password)
	if err != nil {
		logger.Warnf("Login failed for user %s: %v", loginRequest.Username, err)
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{
			"errNo": model.ErrorCode.LoginError,
			"msg":   err.Error(),
		})
		return
	}

	// 登录成功，设置会话
	session := sessions.Get(ctx)
	session.Set(UserIDKey, user.ID)
	session.Set(UsernameKey, user.Username)
	session.Set(IsLoggedInKey, true)
	session.Set(LastAccessKey, time.Now().Unix())
	logger.Infof("User %s.%v logged in successfully", user.Username, user.ID)
	// 返回成功的JSON响应
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			UserIDKey:   user.ID,
			UsernameKey: user.Username,
		},
	})
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

func CreateUser(ctx iris.Context) {
	var userRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		Mobile   string `json:"mobile"`
	}
	err := ctx.ReadJSON(&userRequest)
	if err != nil {
		logger.Warnf("Create user failed: %v", err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"errNo":   400,
			"message": "无效的JSON格式: " + err.Error(),
		})
		return
	}

	userService := &UserService{model.DB.DB()}
	user := &model.AdminUser{
		Username: userRequest.Username,
		Password: userRequest.Password,
		Email:    userRequest.Email,
		RealName: userRequest.RealName,
		Mobile:   userRequest.Mobile,
		Status:   true, // 默认启用
	}

	err = userService.CreateUser(user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"errNo":   500,
			"message": "创建用户失败: " + err.Error(),
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model.ErrorCode.SUCCESS,
	})
}
