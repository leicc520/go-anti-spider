package main

import (
	"github.com/gin-gonic/gin"
)

// 请求提交一个url地址数据资料信息
func Router(app *gin.Engine) {
	app.GET("_/s.js", antiJs)
	app.GET("_/icon.jpg", antiImage)
	app.GET("anti/ping", antiPing)
	app.GET("anti/report", antiReport)
}
