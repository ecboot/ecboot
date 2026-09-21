.PHONY: build test lint gen migrate-up migrate-fresh run check-stub

# 宪法 IV 验证命令的快捷形态
build:
	go build ./...

test:
	go test ./...

lint:
	golangci-lint run

# gf gen dao：从数据库生成 internal/dao 与 internal/model（勿手改生成物）
gen:
	gf gen dao

# 迁移：应用到本地库（golang-migrate CLI）
migrate-up:
	migrate -path migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/ecboot" up

# 空库全量重放验证（危险：drop 后重建）
migrate-fresh:
	mysql -h 127.0.0.1 -P 13306 -uroot -pverysecret -e "DROP DATABASE IF EXISTS ecboot_fresh; CREATE DATABASE ecboot_fresh CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
	migrate -path migrations -database "mysql://root:verysecret@tcp(127.0.0.1:13306)/ecboot_fresh" up

run:
	go run main.go

# 进度验收（specs/PROGRESS.md DoD 依据）：统计各渠道 controller 桩数，用数字说话
check-stub:
	@echo "渠道    总端点  未实现桩"; \
	for ch in admin common shop user; do \
		total=`ls internal/controller/$$ch/ 2>/dev/null | grep -c '^.*_v1_.*\.go$$'`; \
		stub=`grep -l "CodeNotImplemented" internal/controller/$$ch/*.go 2>/dev/null | wc -l`; \
		echo "$$ch    $$total    $$stub"; \
	done; \
	echo "未连线桩合计: `grep -l "CodeNotImplemented" internal/controller/*/*.go 2>/dev/null | wc -l` / 210"
