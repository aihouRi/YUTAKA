// backend/internal/handler/middleware/auth_middleware.go
package middleware

import (
	"backend/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware はリクエストヘッダーからJWTを検証するミドルウェアです。
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストヘッダーから `Authorization` を取得します。
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// ヘッダーが存在しない場合はエラー
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証ヘッダーが必要です。"})
			return
		}

		// ヘッダーの形式が "Bearer <token>" であることを確認します。
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証ヘッダーの形式が正しくありません。"})
			return
		}
		
		tokenString := parts[1]

		// トークンを検証します。
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			// トークンが無効な場合（期限切れ、署名不正など）
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効な認証トークンです。", "details": err.Error()})
			return
		}

		// 検証成功。クレームからの情報を Gin のコンテキストに保存します。
		// これにより、後続のハンドラでユーザー情報にアクセスできます。
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 次のミドルウェアまたはハンドラに処理を渡します。
		c.Next()
	}
}