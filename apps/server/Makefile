.PHONY: build test lint gen migrate-up migrate-fresh run

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
