# ECBOOT

电商平台 Monorepo —— 包含 Spring Boot 后端与三个前端应用（Web 商城、管理后台、移动端）。

> 项目当前处于初始脚手架阶段：后端仅 `start/` 模块包含可运行代码，其余 Maven 模块为规划中的分层占位。

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Java 25、Spring Boot 4.1.1（WebMVC、Security、JPA + Flyway、Redis、Elasticsearch、Quartz、Mail、WebSocket 等） |
| Web 商城 | TanStack Start（React 19）、TanStack Router、Tailwind CSS 4、Vite 8 |
| 管理后台 | TanStack Start（React 19）、TanStack Router、Tailwind CSS 4、Vite 8 |
| 移动端 | uni-app（Vue 3），可编译到微信/支付宝等小程序、H5 与原生 App |

## 目录结构

```
├── start/                 # 可运行的 Spring Boot 应用（唯一包含代码的后端模块）
├── apps/                  # 各渠道 API 模块：user / shop / admin / common（占位）
├── services/              # 领域服务模块：user / shop（占位）
├── infrastructure/        # 公共基础库：common / infra-core（占位）
├── dependencies/          # 依赖 BOM（占位）
├── frontend/
│   ├── ecboot-web/        # Web 商城前端
│   ├── ecboot-admin/      # 管理后台前端
│   └── ecboot-mobile/     # 移动端（uni-app）
├── compose.yaml           # 本地基础设施：MySQL、Redis、Elasticsearch
└── scripts/               # 脚本工具
```

## 快速开始

### 环境要求

- JDK 25
- Docker（本地 MySQL / Redis / Elasticsearch）
- Node.js（前端）

### 后端

```bash
cd start
../mvnw spring-boot:run    # 启动应用；Spring Boot docker-compose 支持会自动拉起 compose.yaml 中的基础设施（需 Docker 已启动）
../mvnw test               # 运行测试
../mvnw test -Dtest=EcbootApplicationTests   # 运行单个测试
```

> 注意：根 `pom.xml` 不是聚合器（无 `<modules>`），请在 `start/` 目录内构建。

### 前端

Web 商城 / 管理后台（TanStack Start）：

```bash
cd frontend/ecboot-web     # 或 frontend/ecboot-admin
npm install
npm run dev                # 开发服务器端口 3000（两个项目相同，同时运行需修改其一）
npm run build              # 生产构建
```

移动端（uni-app）：

```bash
cd frontend/ecboot-mobile
npm install
npm run dev:h5             # H5 开发
npm run dev:mp-weixin      # 微信小程序开发构建
npm run build:mp-weixin    # 微信小程序生产构建
```

## 开发规范

- 提交信息使用中文 Conventional Commits 格式，如 `feat: 新增用户登录`、`chore: 升级依赖`。
- 使用 GitHub Spec Kit 进行规范化开发（见 `.specify/`，配合 `speckit-*` 技能）。
