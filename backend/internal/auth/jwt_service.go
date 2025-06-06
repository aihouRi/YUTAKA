// backend/internal/auth/jwt_service.go
package auth

import (
	"backend/internal/config" // 設定情報を取得するため
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Claims はJWTに含まれるカスタムクレームを定義します。
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken はユーザーIDとユーザー名を受け取り、JWTトークン文字列を生成します。
func GenerateToken(userID int64, username string) (string, error) {
	// トークンの有効期限を設定します。例として24時間とします。
	expirationTime := time.Now().Add(24 * time.Hour)

	// クレームを作成します。
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 有効期限 (ExpiresAt) はUnixタイムスタンプで指定します。
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// 発行日時 (IssuedAt)
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 発行者 (Issuer) - オプション
			Issuer: "YUTAKA",
		},
	}

	// ヘッダーとペイロード（クレーム）を含むトークンオブジェクトを作成します。
	// 署名アルゴリズムとして HS256 を使用します。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 設定ファイルからJWTの秘密鍵を取得します。
	jwtKey := []byte(config.AppConfig.JWTSecretKey)

	// トークンに署名し、完全なトークン文字列を取得します。
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}