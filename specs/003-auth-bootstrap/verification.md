# Verification Record: 003 认证引导纵切片（宪法 IV 留证）

日期：2026-09-20 | 环境：Go 1.25 / GoFrame v2.10.3 / compose（MySQL 8.4@13306、Redis@6379）

## 自动化证据

- `go build ./...` ✓ | `go test ./...` ✓（TestAuthFlow 端到端 PASS，0.34s）
- **端到端序列**（TestAuthFlow）：图形码一次性（同凭证复用拒）→ 发短信码（mock）→ 注册即登录（isNew=true, token+refreshToken 发放）→ **库内零明文**（phone=密文、phone_hash 非空）→ 会话 Validate 有效 → 微信归并同手机号=同一账号 → 另一微信身份同号=20004 冲突 → 重复登录幂等 isNew=false → **休眠拨表 91 天**：微信静默被拒、短信通道放行
- 运行时：`/common/captcha` 真实输出 `{code:0, data:{captchaKey, captchaImg:"data:image/png;base64,..."}}`；路由表 219 端点

## 修复记录（实现过程）

- GoFrame contrib 驱动注册缺失（mysql/nosql redis）→ `bootstrap/driver.go` 空导入 + cmd.go 接线 bootstrap
- phone key 开发默认值须恰好 32 字节（AES-256）
- gvalid.Error 为接口、gerror.Error.Code() 返回 gcode.Code——测试提取器适配两层
