// ctx.go 会话主体 ctx 键（007-admin-base: 从 middleware 下沉至 consts——
// service 层自守卫需读取操作者 ID, 直接 import middleware 会与 middleware→service 依赖成环）。
// 键使用自定义类型（staticcheck SA1029: 避免与原生 string 键碰撞）。
package consts

// CtxKey ctx 键类型（防原生 string 键碰撞）。
type CtxKey string

// 当前会话主体 ID 键（会员或管理员, 随渠道隔离, 见 research D1）。
const CtxUserId CtxKey = "currentUserId"

// 当前访问凭证键（登出销毁用）。
const CtxToken CtxKey = "currentToken"
