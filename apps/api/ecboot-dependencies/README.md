# ecboot-dependencies（版本仲裁 BOM）

全仓库**第三方/框架版本的唯一声明点**——独立 BOM（不继承 ecboot-parent，避免"根导入 BOM + BOM 继承根"的导入自环）。

## 功能内容

| 项 | 说明 |
|---|---|
| 框架版本字面量 | `spring-boot.version=4.1.1`（全仓库唯一出现处） |
| BOM 导入 | `spring-boot-dependencies`（传递管理全部 Spring 家族与常用第三方版本） |
| 第三方追加位 | 未来任何新第三方依赖（如短信 SDK、OSS SDK）**先在此登记版本**，业务模块零版本引用 |

## 职责边界

- **做**：只做 `dependencyManagement`（import 作用域仅消费版本管理，不传递依赖）
- **不做**：不声明任何 `dependencies`、不依赖任何 `org.juling.ecboot` 业务构件（结构性保证 BOM 纯净）

## 使用契约

1. 模块新增第三方依赖：先在本 BOM 登记版本 → 模块内无版本号引用；
2. 升级框架：仅改 `spring-boot.version` 一处，全仓库生效；
3. 根 POM（ecboot-parent）通过 `dependencyManagement` import 本 BOM，reactor 内解析。

相关契约：`specs/001-maven-module-deps/contracts/module-contracts.md` §BOM。
