// Package captcha 图形验证码组件（infra 技术组件, 不与任何渠道耦合）。
// 一次性语义由 Redis GETDEL 原子保证；TTL 走 system_config（默认 300s）。
package captcha

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	keyPrefix = "captcha:img:"
	// TTLFallback 配置缺失时的代码内置默认（配置是覆盖层, 宪法 V）。
	TTLFallback = 300 // 秒
	// AnswerLength 验证码字符数。
	AnswerLength = 4
	// answerChars 去混淆字符集（无 0/O/1/I）。
	answerChars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

// Generate 生成图形验证码：返回凭证 key、Base64 PNG（data-URI）与答案。
// 答案仅服务端留存（Redis），不随响应返回。
func Generate(ctx context.Context) (key, imgBase64, answer string, err error) {
	answer = randomAnswer(AnswerLength)
	keyBytes := make([]byte, 16)
	if _, err = rand.Read(keyBytes); err != nil {
		return "", "", "", err
	}
	key = hex.EncodeToString(keyBytes)

	pngBytes, err := draw(answer)
	if err != nil {
		return "", "", "", err
	}

	_, err = g.Redis().Do(ctx, "SET", keyPrefix+key, normalize(answer), "EX", ttlSeconds(ctx))
	if err != nil {
		return "", "", "", err
	}
	imgBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	return key, imgBase64, answer, nil
}

// Verify 一次性校验：GETDEL 原子取删——无论成败凭证立即失效（spec FR-005）。
// 答案忽略大小写；key 不存在（过期/已消费）返回 false。失败语义由调用方映射（20001）。
func Verify(ctx context.Context, key, answer string) (bool, error) {
	if key == "" || answer == "" {
		return false, nil
	}
	v, err := g.Redis().Do(ctx, "GETDEL", keyPrefix+key)
	if err != nil {
		return false, err
	}
	stored := v.String()
	if stored == "" {
		return false, nil
	}
	return stored == normalize(answer), nil
}

func normalize(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func randomAnswer(n int) string {
	out := make([]byte, n)
	for i := range out {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(answerChars))))
		out[i] = answerChars[idx.Int64()]
	}
	return string(out)
}

func ttlSeconds(ctx context.Context) int {
	v, err := g.DB().GetOne(ctx,
		"SELECT value FROM system_config WHERE code=? AND status=1 AND deleted=0", "captcha.image.ttl_seconds")
	if err != nil || v.IsEmpty() {
		return TTLFallback
	}
	if n := v["value"].Int(); n > 0 {
		return n
	}
	return TTLFallback
}
