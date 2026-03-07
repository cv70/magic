package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// 使用简单的 HMAC-SHA256 实现 JWT，避免添加外部依赖
// 格式: header.payload.signature

const jwtSecret = "your-secret-key-change-in-production" // TODO: 应从环境变量读取

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"` // 过期时间戳
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64, email string) string {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // token 有效期 24 小时

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Exp:    expiresAt.Unix(),
	}

	// 创建 header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// 创建 payload
	claimsJSON, _ := json.Marshal(claims)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := hmac.New(sha256.New, []byte(jwtSecret))
	signature.Write([]byte(message))
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))

	return message + "." + signatureEncoded
}

// ParseToken 解析并验证 JWT token
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 分割 token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	// 验证签名
	message := headerEncoded + "." + payloadEncoded
	signature := hmac.New(sha256.New, []byte(jwtSecret))
	signature.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(signature.Sum(nil))

	if signatureEncoded != expectedSignature {
		return nil, errors.New("invalid signature")
	}

	// 解析 payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, errors.New("failed to unmarshal claims")
	}

	// 检查过期时间
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}
