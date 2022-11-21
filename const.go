package go_anti_spider

import (
	"github.com/go-redis/redis"
)

const (
	visitorId   = "_uuid"
	fingerPrint = "_fp"
	visitPath   = "_sp"
	ipFilter    = "ip@filter"
)

var (
	cRedis *redis.Client = nil
)

// 设置redis客户端处理逻辑
func Inject(rds *redis.Client) {
	cRedis = rds
}
