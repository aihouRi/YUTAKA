// backend/internal/config/config.go
package config

import (
	"errors"
	"log"
	"os"
)

// Config はアプリケーションの設定を保持します。
type Config struct {
	ServerPort   string // ":8080"
	DatabaseDSN  string // データベース接続文字列
	JWTSecretKey string // JWT署名用の秘密鍵
}

// AppConfig はロードされた設定を保持するグローバル変数です。
var AppConfig Config

// LoadConfig は環境変数または.envファイルから設定をロードします。
func LoadConfig() error {
	AppConfig.ServerPort = os.Getenv("SERVER_PORT")
	AppConfig.DatabaseDSN = os.Getenv("DATABASE_DSN")
	AppConfig.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	if AppConfig.DatabaseDSN == "" { return errors.New("DATABASE_DSN is not set") }
	if AppConfig.JWTSecretKey == "" { return errors.New("JWT_SECRET_KEY is not set") }

	if AppConfig.ServerPort == "" { AppConfig.ServerPort = ":8080" }

	log.Println("設定が正常にロードされました。")
	return nil
}

