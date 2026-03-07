package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 使用 SHA256 简单哈希密码（生产环境应使用 bcrypt）
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	return hashedPassword == HashPassword(password)
}
