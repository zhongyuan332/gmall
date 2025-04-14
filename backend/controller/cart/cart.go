package cart

import (
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/controller/common"
	"github.com/zhongyuan332/gmall/backend/logger"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"github.com/zhongyuan332/gmall/sessions"
)

// Create 购物车中添加商品
func Create(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	var cart model2.Cart

	if ctx.ReadJSON(&cart) != nil {
		SendErrJSON("参数错误", ctx)
		return
	}

	if cart.Count <= 0 {
		SendErrJSON("count不能小于0", ctx)
		return
	}

	var product model2.Product
	if model2.DB.First(&product, cart.ProductID).Error != nil {
		SendErrJSON("错误的商品id", ctx)
		return
	}
	session := sessions.Get(ctx)

	openID := session.GetString("weAppOpenID")
	logger.Debugf("user add cart openID: %s, productId:%v", openID, cart.ProductID)
	if openID == "" {
		SendErrJSON("登录超时", ctx)
		return
	}

	cart.OpenID = openID
	if model2.DB.Create(&cart).Error != nil {
		SendErrJSON("error", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"id": cart.ID,
		},
	})
	return
}
