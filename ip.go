package go_anti_spider

import (
	"strconv"
	"time"
)

const (
	ipFieldExp = "exp"
	ipFieldNum = "num"
)

// 统计ip的行为数据分析
type ipSt struct {
	ip string
}

func NewIp(ip string) *ipSt {
	return &ipSt{ip: ip}
}

// 根据历史记录返回是否屏蔽该IP
func (s *ipSt) IsFilter() bool {
	if cmd := cRedis.SIsMember(ipFilter, s.ip); cmd == nil || !cmd.Val() { //没有被屏蔽
		return false
	}
	result := func() bool {
		ckey := "ip@" + s.ip
		expStr := cRedis.HGet(ckey, ipFieldExp).Val()
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
		cRedis.SRem(ipFilter, s.ip)
	}
	return result
}

// 设置ip屏蔽的逻辑
func (s *ipSt) SetFilter() {
	cRedis.SAdd(ipFilter, s.ip)
	ckey := "ip@" + s.ip
	numStr := cRedis.HGet(ckey, ipFieldNum).Val()
	num, _ := strconv.Atoi(numStr)
	num += 1
	expire := time.Now().Unix() + int64(num)*int64(time.Minute*10)
	cRedis.MSet(ipFieldExp, expire, ipFieldNum, num)
}
