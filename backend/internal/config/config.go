// backend/internal/config/config.go
package config

import (
	"errors"
	"log"
	"os"
)

// Config はアプリケーションの設定を保持します。
type Config struct {
	ServerPort   string // 例: ":8080"
	DatabaseDSN  string // データベース接続文字列
	JWTSecretKey string // JWT署名用の秘密鍵
}

// AppConfig はロードされた設定を保持するグローバル変数です。
// 注意: グローバル変数の使用は慎重に。より大きなアプリケーションでは依存性注入を検討してください。
var AppConfig Config

// LoadConfig は環境変数または.envファイルから設定をロードします。
func LoadConfig() error {
	// TODO: 各設定値を環境変数から読み込む処理を実装 (os.Getenv)
	AppConfig.ServerPort = os.Getenv("SERVER_PORT")
	AppConfig.DatabaseDSN = os.Getenv("DATABASE_DSN")
	AppConfig.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	// TODO: 必須の設定値が空の場合のエラーハンドリング
	if AppConfig.DatabaseDSN == "" { return errors.New("DATABASE_DSN is not set") }
	if AppConfig.JWTSecretKey == "" { return errors.New("JWT_SECRET_KEY is not set") }

	// TODO: ServerPort にデフォルト値を設定 (もし環境変数で設定されていなければ)
	if AppConfig.ServerPort == "" { AppConfig.ServerPort = ":8080" }

	log.Println("設定が正常にロードされました。")
	return nil
}

