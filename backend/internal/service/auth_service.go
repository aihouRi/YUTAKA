// backend/internal/service/auth_service.go
package service

import (
	"backend/internal/auth"       // パスワードチェック用
	"backend/internal/repository" // ユーザー取得用
	"errors"
	"fmt"
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

func ChangePassword(userID int64, currentPassword string, newPassword string) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("ユーザー情報の取得に失敗しました: %v", err)
	}

	check := auth.CheckPasswordHash(currentPassword, user.Password)

	if !check {
		return errors.New("現在のパスワードが正しくないです。")
	}

	if len(newPassword) < 8 {
		return errors.New("新しいパスワードの形式が正しくないです新しいパスワードは8文字以上である必要があります。")
	}

	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("パスワードのハッシュ化に失敗しました: %v", err)
	}

	_, err = repository.UpdateUserPassword(newHashedPassword, userID)
	if err != nil {
		return fmt.Errorf("パスワードの更新に失敗しました: %v", err)
	}
	return nil
}
