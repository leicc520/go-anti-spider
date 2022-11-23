package main

// 主函数核心业务逻辑

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

const (
	VERSION = "v1.0.1"
)

var (
	gWindow = js.ValueOf(nil)
	gDebug  = false
)

// 写入日志逻辑
func writeLog(args ...interface{}) {
	if gDebug {
		fmt.Println(args...)
	}
}

// 设置打开调试模式
func antiDebug(_ js.Value, args []js.Value) any {
	if len(args) > 0 { //参数大于1的情况
		if s := args[0].Bool(); s {
			gDebug = true
			return nil
		}
	}
	gDebug = false
	return nil
}

// 获取token数据资料信息
func _token(ox bool, data string) string {
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(data+timeStr)))
	preHash, aftHash, preIdx, aftIdx := make([]byte, 4), make([]byte, 4), 0, len(hashStr)-16
	if ox == false { //取计算
		preIdx += 1
		aftIdx += 1
	}
	for i := 0; i < 4; i++ { //截取指定的字符
		preHash[i] = []byte(hashStr)[preIdx+i*2]
		aftHash[i] = []byte(hashStr)[aftIdx+i*2]
	}
	tokenStr := string(preHash) + timeStr + string(aftHash)
	return tokenStr
}

// 开始调试模式的开关
func antiToken(_ js.Value, _ []js.Value) any {
	defer func() {
		if err := recover(); err != nil {
			writeLog(-1, err)
		}
	}()
	agent := gWindow.Get("navigator").Get("userAgent").String()
	screen := gWindow.Get("screen")
	winScaleStr := screen.Get("width").String() + "x" + screen.Get("height").String()
	tokenStr := _token(false, agent+winScaleStr)
	gmtStr := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	jsSetCookie := `document.cookie = "_at="` + tokenStr + `";expires="` + gmtStr + `";path=/";`
	gWindow.Call("eval", jsSetCookie) //更新处理逻辑
	return nil
}

// 数据加密逻辑
func antiCrypt(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return nil
	}
	dataStr := args[0].String()
	return _token(true, dataStr)
}

// 解析器主函数
func main() {
	gWindow = js.Global()
	gWindow.Set("_av", VERSION)
	gWindow.Set("_at", js.FuncOf(antiToken))
	gWindow.Set("_en", js.FuncOf(antiCrypt))
	gWindow.Set("_dg", js.FuncOf(antiDebug))
	//自动更新token的处理逻辑 30秒刷新token逻辑
	jsAutoReSet := `window.setInterval(function(){window._at()}, 300000)`
	gWindow.Call("eval", jsAutoReSet)
	select {}
}
