// backend/internal/auth/jwt_service.go
package auth

import (
	"backend/internal/config" // 設定情報を取得するため
	"errors"
	"fmt"
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

// ValidateToken はJWTトークン文字列を検証し、有効であればクレームを返します。
func ValidateToken(tokenString string) (*Claims, error) {
	// 設定からJWTの秘密鍵を取得します。
	jwtKey := []byte(config.AppConfig.JWTSecretKey)

	// カスタムクレームを使ってトークンを解析します。
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムが期待通り（HS256）であることを確認します。
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("予期しない署名アルゴリズムです: %v", token.Header["alg"])
		}
		// 検証用のキーを返します。
		return jwtKey, nil
	})

	if err != nil {
		// パース中にエラーが発生した場合（例：署名が不正、有効期限切れなど）
		return nil, fmt.Errorf("無効なトークンです: %w", err)
	}

	// トークンからクレームを抽出し、その有効性を確認します。
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("無効なトークンです")
}