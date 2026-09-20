// driver.go 数据库/Redis 驱动注册（GoFrame v2 贡献包需显式导入注册适配器）。
package bootstrap

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)
