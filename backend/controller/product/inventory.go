package product

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/zhongyuan332/gmall/backend/controller/common"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"strconv"
)

func combinationInventory(productID uint, properties []model2.Property) []model2.Inventory {
	var inventories []model2.Inventory
	if len(properties) == 1 {
		for i := 0; i < len(properties[0].PropertyValues); i++ {
			var inventory = model2.Inventory{
				ProductID: productID,
				PropertyValues: []model2.PropertyValue{
					properties[0].PropertyValues[i],
				},
			}
			inventories = append(inventories, inventory)
		}
	} else {
		theInventories := combinationInventory(productID, properties[1:])
		property := properties[0]
		for i := len(theInventories) - 1; i >= 0; i-- {
			for j := 0; j < len(property.PropertyValues); j++ {
				var inventory = model2.Inventory{
					ProductID:      productID,
					PropertyValues: theInventories[i].PropertyValues,
				}
				inventory.PropertyValues = append(inventory.PropertyValues, property.PropertyValues[j])
				inventories = append(inventories, inventory)
			}
			theInventories = append(theInventories[:i], theInventories[i+1:]...)
		}
	}
	return inventories
}

// SaveInventory 更新商品库存
func SaveInventory(ctx iris.Context) {
	SendErrJSON := common.SendErrJSON
	type InventoryData struct {
		ID    uint `json:"id"`
		Count uint `json:"count"`
	}
	type RequestData struct {
		ProductID   uint            `json:"productID"`
		Inventories []InventoryData `json:"inventories"`
	}

	var reqData RequestData
	var product model2.Product

	if err := ctx.ReadJSON(&reqData); err != nil {
		fmt.Println(err.Error())
		SendErrJSON("参数无效", ctx)
		return
	}

	if err := model2.DB.First(&product, reqData.ProductID).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("错误的商品id", ctx)
		return
	}

	if err := model2.DB.Model(&product).Related(&product.Inventories).Error; err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	for i := 0; i < len(reqData.Inventories); i++ {
		found := false
		for j := 0; j < len(product.Inventories); j++ {
			if reqData.Inventories[i].ID == product.Inventories[j].ID {
				found = true
				break
			}
		}
		if !found {
			SendErrJSON("无效的库存id("+strconv.Itoa(int(reqData.Inventories[i].ID))+")", ctx)
			return
		}
	}

	tx := model2.DB.Begin()
	var count uint
	for i := 0; i < len(reqData.Inventories); i++ {
		err := tx.Exec("UPDATE inventories SET count = ? WHERE id = ?", reqData.Inventories[i].Count, reqData.Inventories[i].ID).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			SendErrJSON("error", ctx)
			return
		}
		count += reqData.Inventories[i].Count
	}

	if err := tx.Exec("UPDATE products SET total_inventory = ? WHERE id = ?", count, reqData.ProductID).Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	tx.Commit()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"errNo": model2.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  iris.Map{},
	})
}
