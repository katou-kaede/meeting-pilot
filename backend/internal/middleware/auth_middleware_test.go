package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"meeting-pilot/internal/auth"
)

func TestRequireAuth(t *testing.T) {
	t.Setenv(
		"JWT_SECRET",
		"test-secret-key-for-auth-middleware",
	)

	validToken, err := auth.GenerateToken(10)
	if err != nil {
		t.Fatalf(
			"テスト用JWTの生成に失敗しました: %v",
			err,
		)
	}

	tests := []struct {
		name           string
		cookie         *http.Cookie
		wantStatusCode int
		wantNextCalled bool
		wantUserID     int64
	}{
		{
			name: "有効なCookieなら認証成功",
			cookie: &http.Cookie{
				Name:  "access_token",
				Value: validToken,
			},
			wantStatusCode: http.StatusOK,
			wantNextCalled: true,
			wantUserID:     10,
		},
		{
			name:           "Cookieがなければ認証失敗",
			cookie:         nil,
			wantStatusCode: http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name: "不正なJWTなら認証失敗",
			cookie: &http.Cookie{
				Name:  "access_token",
				Value: "invalid-token",
			},
			wantStatusCode: http.StatusUnauthorized,
			wantNextCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			nextCalled := false

			next := func(c echo.Context) error {
				nextCalled = true

				userID, ok := c.Get(
					UserIDContextKey,
				).(int64)

				if !ok {
					t.Error(
						"ContextからUserIDを取得できません",
					)
				}

				if userID != tt.wantUserID {
					t.Errorf(
						"UserID = %d, want %d",
						userID,
						tt.wantUserID,
					)
				}

				return c.NoContent(http.StatusOK)
			}

			handler := RequireAuth()(next)

			if err := handler(c); err != nil {
				t.Fatalf(
					"Handler実行時にエラーが発生しました: %v",
					err,
				)
			}

			if rec.Code != tt.wantStatusCode {
				t.Errorf(
					"StatusCode = %d, want %d",
					rec.Code,
					tt.wantStatusCode,
				)
			}

			if nextCalled != tt.wantNextCalled {
				t.Errorf(
					"nextCalled = %v, want %v",
					nextCalled,
					tt.wantNextCalled,
				)
			}
		})
	}
}