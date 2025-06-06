// backend/internal/auth/password_service.go
package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword は平文のパスワードを受け取り、bcryptハッシュを生成します。
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("パスワードのハッシュ化に失敗しました: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash は平文のパスワードとハッシュ化されたパスワードを比較します。
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // エラーがなければパスワードは一致 (如果没有错误则密码一致)
}