// backend/internal/service/auth_service.go
package service

import (
	"errors"
	"fmt"
	"backend/internal/auth"      // パスワードチェック用
	"backend/internal/repository" // ユーザー取得用
)

// Login はユーザー名とパスワードを受け取り、認証を試みます。
// 成功した場合はJWTトークンを、失敗した場合はエラーを返します。
func Login(username string, password string) (string, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("service.Login: ユーザー情報の取得中にエラーが発生しました: %w", err)
	}

	// ユーザーが存在するかどうかを確認します。
	if user == nil {
		return "", errors.New("ユーザー名またはパスワードが正しくありません")
	}

	// パスワードが正しいかを確認します。
	passwordIsValid := auth.CheckPasswordHash(password, user.Password)
	if !passwordIsValid {
		return "", errors.New("ユーザー名またはパスワードが正しくありません")
	}

	// パスワードが正しい場合、JWTを生成します。
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
	    return "", fmt.Errorf("service.Login: トークンの生成に失敗しました: %w", err)
	}
	return token, nil
}