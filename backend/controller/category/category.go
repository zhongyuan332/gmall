package category

import (
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/controller/common"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"github.com/zhongyuan332/gmall/config"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Save 保存分类（创建或更新）
func Save(ctx iris.Context, isEdit bool) {
	SendErrJSON := common.SendErrJSON
	// name, parentId, status, order 必须传的参数
	// remark 非必须
	minOrder := config.ServerConfig.MinOrder
	maxOrder := config.ServerConfig.MaxOrder
	var category model2.Category
	err := ctx.ReadJSON(&category)

	if err != nil {
		SendErrJSON("参数无效", ctx)
		return
	}

	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		SendErrJSON("分类名称不能为空", ctx)
		return
	}

	if utf8.RuneCountInString(category.Name) > config.ServerConfig.MaxNameLen {
		msg := "分类名称不能超过" + strconv.Itoa(config.ServerConfig.MaxNameLen) + "个字符"
		SendErrJSON(msg, ctx)
		return
	}

	if category.Status != model2.CategoryStatusOpen && category.Status != model2.CategoryStatusClose {
		SendErrJSON("status无效", ctx)
		return
	}

	if category.Sequence < minOrder || category.Sequence > maxOrder {
		msg := "分类的排序要在" + strconv.Itoa(minOrder) + "到" + strconv.Itoa(maxOrder) + "之间"
		SendErrJSON(msg, ctx)
		return
	}

	if category.Remark != "" && utf8.RuneCountInString(category.Remark) > config.ServerConfig.MaxRemarkLen {
		msg := "备注不能超过" + strconv.Itoa(config.ServerConfig.MaxRemarkLen) + "个字符"
		SendErrJSON(msg, ctx)
		return
	}

	if category.ParentID != 0 {
		var parentCate model2.Category
		if err := model2.DB.First(&parentCate, category.ParentID).Error; err != nil {
			SendErrJSON("无效的父分类", ctx)
			return
		}
	}

	var updatedCategory model2.Category
	if !isEdit {
		//创建分类
		if err := model2.DB.Create(&category).Error; err != nil {
			SendErrJSON("error", ctx)
			return
		}
	} else {
		//更新分类
		if err := model2.DB.First(&updatedCategory, category.ID).Error; err == nil {
			updatedCategory.Name = category.Name
			updatedCategory.Sequence = category.Sequence
			updatedCategory.ParentID = category.ParentID
			updatedCategory.Status = category.Status
			updatedCategory.Remark = category.Remark
			if err := model2.DB.Save(&updatedCategory).Error; err != nil {
				SendErrJSON("error", ctx)
				return
			}
		} else {
			SendErrJSON("无效的分类id", ctx)
			return
		}
	}

	var categoryJSON model2.Category
	if isEdit {
		categoryJSON = updatedCategory
	} else {
		categoryJSON = category
	}
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"category": categoryJSON,
		},
	})
	return
}

// Create 创建分类
func Create(ctx iris.Context) {
	Save(ctx, false)
}

// Update 更新分类
func Update(ctx iris.Context) {
	Save(ctx, true)
}

// Info 获取分类信息
func Info(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		SendErrJSON("错误的分类id", ctx)
		return
	}

	var category model2.Category
	queryErr := model2.DB.First(&category, id).Error

	if queryErr != nil {
		SendErrJSON("错误的分类id", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"category": category,
		},
	})
}

// AllList 所有的分类列表
func AllList(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	var categories []model2.Category
	pageNo, err := strconv.Atoi(ctx.FormValue("pageNo"))

	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	//默认按创建时间，降序来排序
	var orderStr = "created_at"
	if ctx.FormValue("asc") == "1" {
		orderStr += " asc"
	} else {
		orderStr += " desc"
	}

	offset := (pageNo - 1) * config.ServerConfig.PageSize
	queryErr := model2.DB.Offset(offset).Limit(config.ServerConfig.PageSize).Order(orderStr).Find(&categories).Error

	if queryErr != nil {
		SendErrJSON("error.", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"categories": categories,
		},
	})
}

// List 公开的分类列表
func List(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	var categories []model2.Category

	if model2.DB.Where("status = 1").Order("sequence asc").Find(&categories).Error != nil {
		SendErrJSON("error", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"categories": categories,
		},
	})
}

// UpdateStatus 开启或关闭分类
func UpdateStatus(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	var category model2.Category
	err := ctx.ReadJSON(&category)

	if err != nil {
		SendErrJSON("无效的id或status", ctx)
		return
	}

	id := category.ID
	status := category.Status

	if status != model2.CategoryStatusOpen && status != model2.CategoryStatusClose {
		SendErrJSON("无效的status!", ctx)
		return
	}

	var cate model2.Category
	dbErr := model2.DB.First(&cate, id).Error

	if dbErr != nil {
		SendErrJSON("无效的id!", ctx)
		return
	}

	cate.Status = status

	saveErr := model2.DB.Save(&cate).Error
	if saveErr != nil {
		SendErrJSON("分类状态更新失败", ctx)
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"id":     id,
			"status": status,
		},
	})
}
