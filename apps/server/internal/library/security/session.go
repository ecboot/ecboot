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
//   - token（访问凭证, 短效, 滑动续期）→ session:{token} = userId
//   - refreshToken（刷新凭证, 长效）   → session:refresh:{token} = userId
//
// 语义：refresh 换新双凭证且旧双凭证全失效；登出双凭证同失效（spec FR-015~017）。
type SessionManager struct {
	ttlDays int // 访问凭证天数（system_config: session.ttl_days, 默认 7）
}

func NewSessionManager(ttlDays int) *SessionManager {
	if ttlDays <= 0 {
		ttlDays = 7
	}
	return &SessionManager{ttlDays: ttlDays}
}

// NewSessionManagerFromConfig 会话管理器（TTL 读 system_config: session.ttl_days, 缺失回退 7）。
// 统一入口——避免调用点硬编码 TTL 造成配置漂移（评审 I8）。
func NewSessionManagerFromConfig(ctx context.Context) *SessionManager {
	v, err := g.DB().GetOne(ctx,
		"SELECT value FROM system_config WHERE code='session.ttl_days' AND status=1 AND deleted=0")
	days := 7
	if err == nil && !v.IsEmpty() {
		if n := v["value"].Int(); n > 0 {
			days = n
		}
	}
	return NewSessionManager(days)
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
	if _, err = g.Redis().Do(ctx, "SET", tokenKeyPrefix+token, userId, "EX", int(ttl.Seconds())); err != nil {
		return "", "", err
	}
	// 刷新凭证长效 = 访问凭证 4 倍（30 天口径）
	if _, err = g.Redis().Do(ctx, "SET", refreshKeyPrefix+refreshToken, userId, "EX", int((ttl * 4).Seconds())); err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

// Validate 校验访问凭证：有效则滑动续期并返回 userId。
func (m *SessionManager) Validate(ctx context.Context, token string) (int64, bool, error) {
	if token == "" {
		return 0, false, nil
	}
	v, err := g.Redis().Do(ctx, "GET", tokenKeyPrefix+token)
	if err != nil {
		return 0, false, err
	}
	if v == nil || v.String() == "" {
		return 0, false, nil
	}
	userId := v.Int64()
	// 滑动续期
	ttl := time.Duration(m.ttlDays) * 24 * time.Hour
	_, _ = g.Redis().Do(ctx, "EXPIRE", tokenKeyPrefix+token, int(ttl.Seconds()))
	return userId, true, nil
}

// Refresh 以刷新凭证换发双凭证（旧凭证全部失效, 防重放）。
func (m *SessionManager) Refresh(ctx context.Context, refreshToken string) (token, newRefreshToken string, userId int64, err error) {
	if refreshToken == "" {
		return "", "", 0, fmt.Errorf("刷新凭证为空")
	}
	key := refreshKeyPrefix + refreshToken
	v, err := g.Redis().Do(ctx, "GETDEL", key)
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
	_, err := g.Redis().Do(ctx, "DEL", tokenKeyPrefix+token, refreshKeyPrefix+refreshToken)
	return err
}
