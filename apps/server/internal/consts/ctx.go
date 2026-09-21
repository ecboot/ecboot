// ctx.go 会话主体 ctx 键（007-admin-base: 从 middleware 下沉至 consts——
// service 层自守卫需读取操作者 ID, 直接 import middleware 会与 middleware→service 依赖成环）。
package consts

// CtxUserId 当前会话主体 ID 的 ctx 键（会员或管理员, 随渠道隔离, 见 research D1）。
const CtxUserId = "currentUserId"

// CtxToken 当前访问凭证的 ctx 键（登出销毁用）。
const CtxToken = "currentToken"
