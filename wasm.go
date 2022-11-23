package go_anti_spider

import (
	"crypto/md5"
	"fmt"
	"git.ziniao.com/webscraper/go-orm/log"
	"strconv"
	"time"
)

const (
	hashLen    = 6
	secElaTime = 20
	maxElaTime = 3
)

var (
	GPKey = "597.846.513.246"
)

type WasmSt struct {
	ox   bool   //截取类别
	data string //待验证的参数
	ak   string //请求的校验token
}

// 创建一个wasm对象处理逻辑校验
func NewWasm(ox bool, ak, data string) *WasmSt {
	return &WasmSt{ox: ox, ak: ak, data: data}
}

// 获取token数据资料信息
func (s *WasmSt) _token(timeStr string) string {
	hashStr := fmt.Sprintf("%x", md5.Sum([]byte(s.data+timeStr+GPKey)))
	log.Write(log.INFO, s.data+timeStr+GPKey, hashStr)
	preHash, aftHash, preIdx, aftIdx := make([]byte, hashLen), make([]byte, hashLen), 0, len(hashStr)-16
	if s.ox == false { //取计算
		preIdx += 1
		aftIdx += 1
	}
	for i := 0; i < hashLen; i++ { //截取指定的字符
		preHash[i] = []byte(hashStr)[preIdx+i*2]
		aftHash[i] = []byte(hashStr)[aftIdx+i*2]
	}
	tokenStr := string(preHash) + timeStr + string(aftHash)
	return tokenStr
}

// 截取token的时间戳逻辑
func (s *WasmSt) getTimeStamp() string {
	if len(s.ak) < hashLen*2+10 {
		return "-"
	}
	return s.ak[hashLen+1 : len(s.ak)-hashLen]
}

// 验证请求的参数是否合法的处理逻辑
func (s *WasmSt) Check() bool {
	timeStr := s.getTimeStamp()
	akStr := s._token(timeStr)
	if akStr != s.ak { //数据复合不一致跳过
		return false
	}
	//请求的时间大于预设的时间阈值，说明的模拟记录的请求
	timeInt, _ := strconv.ParseInt(timeStr, 10, 64)
	if time.Now().Unix() > timeInt+secElaTime*maxElaTime {
		return false
	}
	return true
}
