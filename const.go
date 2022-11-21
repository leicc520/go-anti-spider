package go_anti_spider

import (
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

const (
	visitorId   = "_uuid"
	fingerPrint = "_fp"
	visitPath   = "_sp"
	ipFilter    = "ip@filter"
)

var (
	cRedis   *redis.Client      = nil
	cIpFiles map[int]*TmpFileSt = nil
	nIpFiles                    = 1
)

// 设置redis客户端处理逻辑
func Boostrap(rds *redis.Client, nFiles int) {
	cRedis, nIpFiles = rds, nFiles
	cIpFiles = make(map[int]*TmpFileSt)
	for i := 0; i < nIpFiles; i++ {
		cIpFiles[i] = NewTmpFile(ipTmpFile + "-" + strconv.Itoa(nFiles))
	}
}

// 请求IP的业务数据信息收集
func IPCollect(ip, path string) {
	path = strings.ReplaceAll(path, `/`, "-")
	if len(path) < 1 {
		path = "-"
	}
	idx := time.Now().UnixMilli() % int64(nIpFiles)
	cIpFiles[int(idx)].Write(ip + ";" + path)
}
