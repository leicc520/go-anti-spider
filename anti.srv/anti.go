package main

import (
	"github.com/gin-gonic/gin"
)

// 检测图片的处理逻辑
func antiImage(ctx *gin.Context) {
	//生成 vi --使用agent-time 加密
}

// 检测JS的处理逻辑 -其他业务通过这个js调用
func antiJs(ctx *gin.Context) {
	//输出js代码，动态加载image、验证码、Ping、report处理逻辑
}

// 上报API数据逻辑
func antiReport(ctx *gin.Context) {
	//提报请求的API运行情况逻辑
}

// 上报Ping数据逻辑
func antiPing(ctx *gin.Context) {
	//上报提交参数agent window-w window-h 拆解的vi片段
}

// 检测到异常弹出验证码
func antiCaptcha(ctx *gin.Context) {
	//反爬验证码处理逻辑
}
