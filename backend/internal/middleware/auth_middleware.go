package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"meeting-pilot/internal/auth"
)

const UserIDContextKey = "user_id"

func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]string{
						"error": "ログインが必要です",
					},
				)
			}

			claims, err := auth.ParseToken(cookie.Value)
			if err != nil {
				log.Printf(
					"RequireAuth invalid token err=%v",
					err,
				)

				return c.JSON(
					http.StatusUnauthorized,
					map[string]string{
						"error": "ログイン情報が無効または期限切れです",
					},
				)
			}

			// 後続のHandlerでUserIDを取得できるように保存
			c.Set(UserIDContextKey, claims.UserID)

			return next(c)
		}
	}
}
