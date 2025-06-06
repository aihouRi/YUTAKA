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
