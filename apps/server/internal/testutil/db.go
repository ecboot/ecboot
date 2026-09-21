// db.go 测试基座共享配置——**数据库连接的唯一事实源**（012 评审修复轮统一到应用库 ecboot）。
//
// 背景（2026-09-21 实测的漂移）：此前 7 个测试包的 init() 各自硬编码
// `myuser:secret@tcp(127.0.0.1:13306)/mydatabase`（docker 测试库），而应用配置
// （config/config.yaml）、`make migrate-up`、`gf gen dao` 全部指向 `root:root@tcp(127.0.0.1:3306)/ecboot`。
// 结果是三份事实源并行：迁移打到 A 库、测试跑在 B 库、生成物取自 C 库——
// 实测同一张表在两侧的列集不同（ecboot 停在迁移 34），生成物与两库都不完全一致。
//
// 现统一：**测试与应用共用应用库 `ecboot`**（迁移、生成物、测试、应用四者同源）。
// 需要隔离环境（CI/并行分片）时用环境变量覆盖，不要在测试文件里另写连接串。
package testutil

import "os"

// DefaultDSN 默认测试库 = 应用库 ecboot（本机 MySQL 3306）。
//
// 时区口径（009/012 评审的横切口径，必须保留）: 库会话时钟与驱动解释一律 UTC，
// 与 `main.go`/测试基座的 `time.Local = time.UTC` 对齐。
//
// 背景: 本机 MySQL 的 `time_zone=SYSTEM` 是 **+08**，而进程按 UTC 写 DATETIME 墙钟；
// 若不在 DSN 上钉住会话时区，SQL 端 `NOW()` 比进程写入值快 8 小时 ——
// 实测后果: 写"开始时间=+60 分钟"的轮播被 `< NOW()` 判为已开始（在投过滤失效）。
// 历史: batch 03 记录过同族缺陷（读回早 8h、回显再提交每轮漂 8h），修法是"进程与库会话时钟同时锁 UTC"。
//
// 写法约束（2026-09-22 探针实测，勿改）:
//   - 必须用**命名时区** `time_zone=UTC`。GoFrame 会把 DSN 参数里的 `+` 变成空格，
//     故偏移量写法 `'+00:00'` 一律开不了连接（`unknown time zone ' 00:00'`）。
//   - 命名时区要求服务器载入 mysql 库的时区表；本机已用
//     `mysql_tzinfo_to_sql /usr/share/zoneinfo | mysql ... mysql` 载入（此前为空表 → `SET time_zone='UTC'` 报 1298）。
//   - GoFrame 把 `time_zone` 解析进它自己的 Loc 字段并回写 `loc=`，因此一个参数同时管住
//     "驱动如何解释 DATETIME(loc)" 与 "库会话时钟(会话变量)" 两侧 —— 这正是我们要的双侧一致。
const DefaultDSN = "mysql:root:root@tcp(127.0.0.1:3306)/ecboot?loc=UTC&time_zone=UTC"

// RedisAddr 测试用 Redis（本机 6379；测试只用到验证码/会话键，不需隔离库）。
const RedisAddr = "127.0.0.1:6379"

// DSN 测试库连接串（ECBOOT_TEST_DSN 可覆盖）。
func DSN() string {
	if v := os.Getenv("ECBOOT_TEST_DSN"); v != "" {
		return v
	}
	return DefaultDSN
}
