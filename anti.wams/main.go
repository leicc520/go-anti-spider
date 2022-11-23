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
	GPKey   = "x.x.x.x"
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
	defer func() {
		if err := recover(); err != nil {
			writeLog(err)
		}
	}()
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
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(data+timeStr+GPKey)))
	writeLog("D:", data+timeStr+GPKey)
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
			writeLog(err)
		}
	}()
	agent := gWindow.Get("navigator").Get("userAgent").String()
	screen := gWindow.Get("screen")
	winScaleStr := fmt.Sprintf("%dx%d", screen.Get("width").Int(), screen.Get("height").Int())
	tokenStr := _token(false, agent+winScaleStr)
	writeLog("T:", tokenStr)
	gmtStr := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	doc := gWindow.Get("document")
	doc.Set("cookie", js.ValueOf(`_at=`+tokenStr+`;expires=`+gmtStr+`;path=/;`))
	writeLog("J:", `_at=`+tokenStr+`;expires=`+gmtStr+`;path=/;`)
	return nil
}

// 数据加密逻辑
func antiCrypt(_ js.Value, args []js.Value) any {
	defer func() {
		if err := recover(); err != nil {
			writeLog(err)
		}
	}()
	if len(args) < 1 {
		return nil
	}
	dataStr := args[0].String()
	return _token(true, dataStr)
}

// 解析器主函数
func main() {
	gWindow = js.Global()
	gWindow.Set("_v", VERSION)
	gWindow.Set("_a", js.FuncOf(antiToken))
	gWindow.Set("_e", js.FuncOf(antiCrypt))
	gWindow.Set("_d", js.FuncOf(antiDebug))
	//自动更新token的处理逻辑 30秒刷新token逻辑
	jsAutoReSet := `window.setInterval(function(){window._a();}, 300000)`
	jsAutoReSet = `window.setInterval(function(){console.log("1111");}, 3000)`
	gWindow.Call("eval", jsAutoReSet)
	select {}
}
