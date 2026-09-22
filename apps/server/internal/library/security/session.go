package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// SessionManager 不透明令牌会话管理（Redis 承载）：
//   - token（访问凭证, 短效, 滑动续期）→ session:{aud}:{token} = userId
//   - refreshToken（刷新凭证, 长效）   → session:refresh:{aud}:{token} = userId
//
// aud（渠道）∈ user / admin：会话按渠道隔离——修复 admin/user 两表主键同号时
// 凭证跨渠道冒充的越权缺陷（007-admin-base research D1）。
// 语义：refresh 换新双凭证且旧双凭证全失效；登出双凭证同失效（spec FR-015~017）。
type SessionManager struct {
	aud     string // 会话渠道（user/admin）
	ttlDays int    // 访问凭证天数（system_config: session.ttl_days, 默认 7）
}

// NewSessionManager 指定渠道的会话管理器（aud 空值回退 user 以保既有语义）。
func NewSessionManager(aud string, ttlDays int) *SessionManager {
	if aud == "" {
		aud = "user"
	}
	if ttlDays <= 0 {
		ttlDays = 7
	}
	return &SessionManager{aud: aud, ttlDays: ttlDays}
}

// NewSessionManagerFromConfig 会话管理器（TTL 读 system_config: session.ttl_days, 缺失回退 7）。
// 统一入口——避免调用点硬编码 TTL 造成配置漂移（评审 I8）。
func NewSessionManagerFromConfig(ctx context.Context, aud string) *SessionManager {
	v, err := g.DB().GetOne(ctx,
		"SELECT value FROM system_config WHERE code='session.ttl_days' AND status=1 AND deleted=0")
	days := 7
	if err == nil && !v.IsEmpty() {
		if n := v["value"].Int(); n > 0 {
			days = n
		}
	}
	return NewSessionManager(aud, days)
}

const (
	tokenKeyPrefix   = "session:"
	refreshKeyPrefix = "session:refresh:"
)

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Create 登录成功：发放双凭证。
func (m *SessionManager) Create(ctx context.Context, userId int64) (token, refreshToken string, err error) {
	if token, err = newToken(); err != nil {
		return "", "", err
	}
	if refreshToken, err = newToken(); err != nil {
		return "", "", err
	}
	ttl := time.Duration(m.ttlDays) * 24 * time.Hour
	if _, err = g.Redis().Do(ctx, "SET", m.tokenKey(token), userId, "EX", int(ttl.Seconds())); err != nil {
		return "", "", err
	}
	// 刷新凭证长效 = 访问凭证 4 倍（30 天口径）
	if _, err = g.Redis().Do(ctx, "SET", m.refreshKey(refreshToken), userId, "EX", int((ttl * 4).Seconds())); err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

// Validate 校验访问凭证：有效则滑动续期并返回 userId。
func (m *SessionManager) Validate(ctx context.Context, token string) (int64, bool, error) {
	if token == "" {
		return 0, false, nil
	}
	v, err := g.Redis().Do(ctx, "GET", m.tokenKey(token))
	if err != nil {
		return 0, false, err
	}
	if v == nil || v.String() == "" {
		return 0, false, nil
	}
	userId := v.Int64()
	// 滑动续期
	ttl := time.Duration(m.ttlDays) * 24 * time.Hour
	_, _ = g.Redis().Do(ctx, "EXPIRE", m.tokenKey(token), int(ttl.Seconds()))
	return userId, true, nil
}

// PeekUserId 非消费式解析刷新凭证对应的 userId（**不移除**该 refresh token）。
// 用途（N3 收口）: 刷新前校验账号态——禁用/软删会员不得续期。
// 原刷新路径不校验账号态, 且 refresh 每次轮换长寿 4×TTL → 禁用可被近乎无限绕过。
func (m *SessionManager) PeekUserId(ctx context.Context, refreshToken string) (int64, error) {
	if refreshToken == "" {
		return 0, fmt.Errorf("刷新凭证为空")
	}
	v, err := g.Redis().Do(ctx, "GET", m.refreshKey(refreshToken))
	if err != nil {
		return 0, err
	}
	if v == nil || v.String() == "" {
		return 0, fmt.Errorf("刷新凭证失效")
	}
	return v.Int64(), nil
}

// Refresh 以刷新凭证换发双凭证（旧凭证全部失效, 防重放）。
func (m *SessionManager) Refresh(ctx context.Context, refreshToken string) (token, newRefreshToken string, userId int64, err error) {
	if refreshToken == "" {
		return "", "", 0, fmt.Errorf("刷新凭证为空")
	}
	v, err := g.Redis().Do(ctx, "GETDEL", m.refreshKey(refreshToken))
	if err != nil {
		return "", "", 0, err
	}
	if v == nil || v.String() == "" {
		return "", "", 0, fmt.Errorf("刷新凭证失效")
	}
	userId = v.Int64()
	token, newRefresh, err2 := m.Create(ctx, userId)
	if err2 != nil {
		return "", "", 0, err2
	}
	return token, newRefresh, userId, nil
}

// Destroy 登出：访问与刷新凭证同时失效。
func (m *SessionManager) Destroy(ctx context.Context, token, refreshToken string) error {
	_, err := g.Redis().Do(ctx, "DEL", m.tokenKey(token), m.refreshKey(refreshToken))
	return err
}

func (m *SessionManager) tokenKey(token string) string {
	return tokenKeyPrefix + m.aud + ":" + token
}

func (m *SessionManager) refreshKey(refreshToken string) string {
	return refreshKeyPrefix + m.aud + ":" + refreshToken
}
