package visit

import (
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/controller/common"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"strconv"
	"time"
)

// Latest30Day 近30天，每天的PV
func Latest30Day(ctx iris.Context) {
	var userVisit model2.UserVisit
	result := userVisit.Latest30DayPV()
	var data iris.Map
	if result == nil {
		data = iris.Map{
			"list": [0]int{},
		}
	} else {
		data = iris.Map{
			"list": result,
		}
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}

// PV 增加一次页面访问
func PV(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	var err error
	var msg = ""
	var userVisit model2.UserVisit
	userVisit.ClientID = ctx.FormValue("clientId")
	userVisit.OSName = ctx.FormValue("osName")
	userVisit.OSVersion = ctx.FormValue("osVersion")
	userVisit.Language = ctx.FormValue("language")
	userVisit.Country = ctx.FormValue("country")
	userVisit.DeviceModel = ctx.FormValue("deviceModel")
	userVisit.DeviceWidth, err = strconv.Atoi(ctx.FormValue("deviceWidth"))
	if err != nil {
		SendErrJSON("无效的deviceWidth", ctx)
		return
	}
	userVisit.DeviceHeight, err = strconv.Atoi(ctx.FormValue("deviceHeight"))
	if err != nil {
		SendErrJSON("无效的deviceHeight", ctx)
		return
	}
	userVisit.IP = ctx.RemoteAddr()
	userVisit.VisitTime = time.Now()
	userVisit.Referrer = ctx.FormValue("referrer")
	userVisit.URL = ctx.FormValue("url")
	userVisit.BrowserName = ctx.FormValue("browserName")
	userVisit.BrowserVersion = ctx.FormValue("browserVersion")

	if userVisit.ClientID == "" {
		SendErrJSON("clientId不能为空", ctx)
		return
	}

	if msg != "" {
		SendErrJSON("error", ctx)
		return
	}

	if err := model2.DB.Create(&userVisit).Error; err != nil {
		SendErrJSON("error.", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)

	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  iris.Map{},
	})
}
