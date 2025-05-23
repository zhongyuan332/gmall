package user

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/controller/common"
	"github.com/zhongyuan332/gmall/backend/logger"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"github.com/zhongyuan332/gmall/config"
	"github.com/zhongyuan332/gmall/sessions"
	"github.com/zhongyuan332/gmall/utils"
	"net/http"
	"strings"
	"time"
)

// WeAppLogin 微信小程序登录
func WeAppLogin(ctx iris.Context) {
	logger.Debugf("weapp login start...")
	SendErrJSON := common.SendErrJSON
	code := ctx.FormValue("code")
	if code == "" {
		SendErrJSON("code不能为空", ctx)
		return
	}
	appID := config.WeAppConfig.AppID
	secret := config.WeAppConfig.Secret
	CodeToSessURL := config.WeAppConfig.CodeToSessURL
	CodeToSessURL = strings.Replace(CodeToSessURL, "{appid}", appID, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{code}", code, -1)

	resp, err := http.Get(CodeToSessURL)
	if err != nil {
		logger.Errorf("user login error, err:%v", err.Error())
		SendErrJSON("error", ctx)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		SendErrJSON("error", ctx)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	if _, ok := data["session_key"]; !ok {
		fmt.Println("session_key 不存在")
		fmt.Println(data)
		SendErrJSON("error", ctx)
		return
	}

	var openID string
	var sessionKey string
	openID = data["openid"].(string)
	sessionKey = data["session_key"].(string)
	session := sessions.Get(ctx)
	session.Set("weAppOpenID", openID)
	session.Set("weAppSessionKey", sessionKey)

	resData := iris.Map{}
	resData[config.ServerConfig.SessionID] = session.ID()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
}

// SetWeAppUserInfo 设置小程序用户加密信息
func SetWeAppUserInfo(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	type EncryptedUser struct {
		EncryptedData string `json:"encryptedData"`
		IV            string `json:"iv"`
	}
	var weAppUser EncryptedUser

	if ctx.ReadJSON(&weAppUser) != nil {
		SendErrJSON("参数错误", ctx)
		return
	}
	session := sessions.Get(ctx)
	sessionKey := session.GetString("weAppSessionKey")
	if sessionKey == "" {
		SendErrJSON("session error", ctx)
		return
	}

	userInfoStr, err := utils.DecodeWeAppUserInfo(weAppUser.EncryptedData, sessionKey, weAppUser.IV)
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	var user model2.WeAppUser
	if err := json.Unmarshal([]byte(userInfoStr), &user); err != nil {
		SendErrJSON("error", ctx)
		return
	}

	session.Set("weAppUser", user)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  iris.Map{},
	})
	return
}

// YesterdayRegisterUser 昨日注册的用户数
func YesterdayRegisterUser(ctx iris.Context) {
	var user model2.User
	count := user.YesterdayRegisterUser()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"count": count,
		},
	})
}

// TodayRegisterUser 今日注册的用户数
func TodayRegisterUser(ctx iris.Context) {
	var user model2.User
	count := user.TodayRegisterUser()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"count": count,
		},
	})
}

// Latest30Day 近30天，每天注册的新用户数
func Latest30Day(ctx iris.Context) {
	var users model2.UserPerDay
	result := users.Latest30Day()
	var data iris.Map
	if result == nil {
		data = iris.Map{
			"users": [0]int{},
		}
	} else {
		data = iris.Map{
			"users": result,
		}
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}

// Analyze 用户分析
func Analyze(ctx iris.Context) {
	var user model2.User
	now := time.Now()
	nowSec := now.Unix()              //秒
	yesterdaySec := nowSec - 24*60*60 //秒
	yesterday := time.Unix(yesterdaySec, 0)

	yesterdayCount := user.PurchaseUserByDate(yesterday)
	todayCount := user.PurchaseUserByDate(now)
	yesterdayRegisterCount := user.YesterdayRegisterUser()
	todayRegisterCount := user.TodayRegisterUser()
	data := iris.Map{
		"todayNewUser":          todayRegisterCount,
		"yesterdayNewUser":      yesterdayRegisterCount,
		"todayPurchaseUser":     todayCount,
		"yesterdayPurchaseUser": yesterdayCount,
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}
