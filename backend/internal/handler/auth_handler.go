// backend/internal/handler/auth_handler.go
package handler

import (
	"backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginRequest はログインAPIのリクエストボディを定義します。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin はログインリクエストを処理します。
func HandleLogin(c *gin.Context) {
	var req LoginRequest

	// リクエストボディを構造体にバインドし、バリデーションします。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ユーザー名とパスワードを入力してください。",
			"details": err.Error(),
		})
		return
	}

	// 認証サービスを呼び出します。
	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(), 
		})
		return
	}

	// 認証成功。JWTトークンを返します。
	c.JSON(http.StatusOK, gin.H{
		"message": "ログインに成功しました。",
		"token":   token,
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// HandleChangePassword は認証済みユーザーのパスワード変更リクエストを処理します。
func HandleChangePassword(c *gin.Context) {
	// 1. リクエストボディのJSONを構造体にバインドし、バリデーションを行います。
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "入力内容が正しくありません。",
			"details": err.Error(),
		})
		return
	}

	// JWTミドルウェアによって設定されたユーザーIDをコンテキストから取得します。
	userID_any, exists := c.Get("userID")
	if !exists {
		// このエラーは通常、ミドルウェアが正しく適用されていない場合にのみ発生します。
		// サーバー側の設定ミスの可能性が高いため、500エラーを返します。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部でユーザー情報が見つかりませんでした。"})
		return
	}

	// c.Get() は any (interface{}) 型を返すため、int64 に型アサーション（変換）します。
	userID, ok := userID_any.(int64)
	if !ok {
		// 型アサーションが失敗した場合も、サーバー側の問題です。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部でユーザーIDの形式が不正です。"})
		return
	}

	// 認証サービスレイヤーの ChangePassword 関数を呼び出します。
	err := service.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		// サービス層から返されたエラー
		// クライアントに起因するエラーなので、400 Bad Request を返すのが適切です。
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// パスワード変更成功のレスポンスを返します。
	c.JSON(http.StatusOK, gin.H{
		"message": "パスワードが正常に更新されました。",
	})
}
