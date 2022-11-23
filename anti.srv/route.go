package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
)

// 请求提交一个url地址数据资料信息
func Router(app *gin.Engine) {
	app.Use(func(ctx *gin.Context) { //针对wasm的加载
		if strings.Contains(ctx.Request.URL.Path, ".wasm") {
			ctx.Header("Cache-Control", "no-cache")
			ctx.Header("content-type", "application/wasm")
		}
		ctx.Next()
	})
	dir, _ := os.Getwd()
	staticDir := filepath.Join(dir, "static")
	app.Static("/static", staticDir)
	app.GET("_/s.js", antiJs)
	app.GET("anti/ping", antiPing)
	app.GET("anti/rp", antiReport)
	app.GET("anti/ca", antiCaptcha)
}
