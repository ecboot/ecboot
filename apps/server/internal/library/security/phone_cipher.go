// Package security 认证相关技术组件：手机号加密、会话管理。
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// PhoneCipher 手机号加密器：密文 = AES-256-GCM（随机 nonce 前置）；
// 检索哈希 = SHA-256(密钥派生盐 + 手机号)——密钥不入库不入日志（spec FR-010/FR-004）。
type PhoneCipher struct {
	aead cipher.AEAD
	salt []byte
}

// NewPhoneCipher 以 Base64 主密钥构造（32 字节解码后使用）。
func NewPhoneCipher(masterKeyBase64 string) (*PhoneCipher, error) {
	key, err := base64.StdEncoding.DecodeString(masterKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("phone key base64 解码失败: %w", err)
	}
	if len(key) != 32 {
		return nil, errors.New("phone key 必须 decode 为 32 字节(AES-256)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 密钥派生盐（HMAC 主密钥与固定域串, 与加密密钥分离）
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("ecboot:phone_hash:v1"))
	return &PhoneCipher{aead: aead, salt: mac.Sum(nil)}, nil
}

// Encrypt 手机号 → 密文（base64(nonce+ciphertext)）。
func (c *PhoneCipher) Encrypt(phone string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(phone), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt 密文 → 手机号（后台改绑等受控场景使用; 常规检索走哈希）。
func (c *PhoneCipher) Decrypt(cipherText string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	if len(raw) < c.aead.NonceSize() {
		return "", errors.New("密文过短")
	}
	nonce, sealed := raw[:c.aead.NonceSize()], raw[c.aead.NonceSize():]
	plain, err := c.aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// Hash 手机号 → 检索哈希（十六进制; 唯一索引与登录精确检索口径）。
func (c *PhoneCipher) Hash(phone string) string {
	mac := hmac.New(sha256.New, c.salt)
	mac.Write([]byte(phone))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
