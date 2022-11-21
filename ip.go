package go_anti_spider

import (
	"strconv"
	"strings"
	"time"
)

const (
	ipFieldExp  = "exp"
	ipFieldNum  = "num"
	ipLimitBase = 8
	ipExpire    = time.Minute * 10
	ipTmpFile   = "/tmp/ip.anti"
)

// 统计ip的行为数据分析
type ipSt struct {
}

// 创建一个IP解析器
func NewIpParser() *ipSt {
	return &ipSt{}
}

// 根据历史记录返回是否屏蔽该IP
func (s *ipSt) IsFilter(ip, path string) bool {
	if cmd := cRedis.SIsMember(ipFilter, ip); cmd == nil || !cmd.Val() { //没有被屏蔽
		return false
	}
	//被屏蔽的话继续验证 过期时效
	result := func() bool {
		expStr := cRedis.HGet("ip@"+ip, ipFieldExp).Val()
		if len(expStr) < 1 { //过期时间为空的情况
			return false
		}
		expire, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil || expire < time.Now().Unix() {
			return false
		}
		return true
	}()
	if !result { //未被屏蔽的情况
		cRedis.SRem(ipFilter, ip)
	}
	//提交到IP收集临时文件
	IPCollect(ip, strings.Trim(path, `/`))
	return result
}

// 分析IP行为问题
func (s *ipSt) Parse(line string) {
	aStr := strings.SplitN(strings.TrimSpace(line), ";", 2)
	if len(aStr) < 2 { //收集字段小于2的情况
		return
	}
	hKey := "path@" + aStr[0]
	cmd := cRedis.HIncrBy(hKey, aStr[1], 1)
	if cmd == nil || cmd.Val() < ipLimitBase {
		return
	}
	//分析请求的地址情况信息
}

// 设置ip屏蔽的逻辑
func (s *ipSt) SetFilter(ip string, expire int64) {
	cRedis.SAdd(ipFilter, ip)
	cKey, num := "ip@"+ip, 1
	if expire < 1 {
		numStr := cRedis.HGet(cKey, ipFieldNum).Val()
		num, _ = strconv.Atoi(numStr)
		num += 1
		expire = time.Now().Unix() + int64(num)*int64(ipExpire)
	}
	cRedis.MSet(ipFieldExp, expire, ipFieldNum, num)
}
