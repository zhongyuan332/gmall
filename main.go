package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	log "github.com/zhongyuan332/gmall/backend/logger"
	model2 "github.com/zhongyuan332/gmall/backend/model"
	"github.com/zhongyuan332/gmall/backend/route"
	"github.com/zhongyuan332/gmall/config"
	"github.com/zhongyuan332/gmall/sessions"
	"os"
	"strconv"
	"time"
)

func init() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if config.DBConfig.SQLLog {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)

	model2.DB = db
}

func main() {
	// 初始化配置
	log.InitLogger(log.DefaultConfig)

	app := iris.New()
	app.Use(iris.Compression)                // 启用 Gzip 压缩
	app.Configure(iris.WithCharset("UTF-8")) // 设置字符集

	if config.ServerConfig.Debug {
		app.Logger().SetLevel("debug")
		app.Use(logger.New())
	}

	// 创建会话管理器
	sessions.Initialize(config.ServerConfig.SessionID, time.Minute*20)

	app.Use(sessions.Manager.Handler())

	// 注册路由
	route.Route(app)
	// 错误处理
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"errNo": model2.ErrorCode.NotFound,
			"msg":   "Not Found",
			"data":  iris.Map{},
		})
	})

	app.OnErrorCode(500, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"errNo": model2.ErrorCode.ERROR,
			"msg":   "error",
			"data":  iris.Map{},
		})
	})

	app.Listen(":" + strconv.Itoa(config.ServerConfig.Port))
}
