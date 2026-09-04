package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"meeting-pilot/internal/auth"
	"meeting-pilot/internal/model"
	"meeting-pilot/internal/repository"
	"meeting-pilot/internal/validator"
)


// ============================================
// ユーザー登録
// ============================================
func CreateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.CreateUserRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"CreateUser bind failed err=%v",
				err,
			)

			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "入力内容の形式が正しくありません",
				},
			)
		}

		// バリデーション
		if err := validator.ValidateCreateUser(req); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		user, err := repository.CreateUser(db, req)
		if err != nil {
			if errors.Is(err, repository.ErrEmailAlreadyExists) {
				return c.JSON(
					http.StatusConflict,
					map[string]string{
						"error": "このメールアドレスは既に登録されています",
					},
				)
			}

			log.Printf(
				"CreateUser failed email=%s err=%v",
				req.Email,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "ユーザーの登録に失敗しました",
				},
			)
		}

		return c.JSON(
			http.StatusCreated,
			user,
		)
	}
}

// ============================================
// ログイン
// ============================================
func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.LoginRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"Login bind failed err=%v",
				err,
			)

			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "入力内容の形式が正しくありません",
				},
			)
		}

		user, err := repository.GetUserByEmail(
			db,
			req.Email,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]string{
						"error": "メールアドレスまたはパスワードが正しくありません",
					},
				)
			}

			log.Printf(
				"Login get user failed email=%s err=%v",
				req.Email,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "ログインに失敗しました",
				},
			)
		}

		// 利用停止ユーザーはログイン不可
		if !user.IsActive {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "このアカウントは現在利用できません",
				},
			)
		}

		// パスワードを照合
		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(req.Password),
		); err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				map[string]string{
					"error": "メールアドレスまたはパスワードが正しくありません",
				},
			)
		}

		// JWTを生成
		token, err := auth.GenerateToken(user.ID)
		if err != nil {
			log.Printf(
				"Login token generation failed user_id=%d err=%v",
				user.ID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "ログインに失敗しました",
				},
			)
		}

		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		sameSite := http.SameSiteLaxMode
		if cookieSecure {
				sameSite = http.SameSiteNoneMode
		}

		// JWTをHttpOnly Cookieへ保存
		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   cookieSecure, // 本番のHTTPS環境ではtrue(そのCookieをHTTPS通信のときだけ送信するか)
			SameSite: sameSite,
			MaxAge:   60 * 60, // 1時間
		})

		return c.JSON(
			http.StatusOK,
			map[string]any{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		)
	}
}

// ============================================
// ログイン後の認証ユーザー情報取得
// ============================================
func GetCurrentUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ログインユーザーIDを取得
		userID, err := getAuthenticatedUserID(c)
		if err != nil {
			return err
		}

		user, err := repository.GetUserByID(
			db,
			userID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(
					http.StatusUnauthorized,
					map[string]string{
						"error": "ユーザーが見つかりません",
					},
				)
			}

			log.Printf(
				"GetCurrentUser failed user_id=%d err=%v",
				userID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "ユーザー情報の取得に失敗しました",
				},
			)
		}

		if !user.IsActive {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "このアカウントは現在利用できません",
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			user,
		)
	}
}

// ============================================
// ログアウト
// ============================================
func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		sameSite := http.SameSiteLaxMode
		if cookieSecure {
				sameSite = http.SameSiteNoneMode
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   cookieSecure, // 本番HTTPSではtrue
			SameSite: sameSite,
			MaxAge:   -1,
		})

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// ユーザー検索
// ============================================
func SearchUsersForMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ログインユーザーIDを取得
		userID, err := getAuthenticatedUserID(c)
		if err != nil {
			return err
		}

		meetingID, err := strconv.ParseInt(
			c.Param("id"),
			10,
			64,
		)
		if err != nil {
			log.Printf(
				"SearchUsersForMeeting invalid meeting_id=%s user_id=%d err=%v",
				c.Param("id"),
				userID,
				err,
			)

			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "会議IDが不正です",
				},
			)
		}

		// ログインユーザーの権限チェック
		_, err = getMeetingUserRole(
			c,
			db,
			meetingID,
			userID,
		)
		if err != nil {
			return err
		}

		keyword := c.QueryParam("keyword")

		users, err := repository.SearchUsersForMeeting(
			db,
			meetingID,
			keyword,
		)
		if err != nil {
			log.Printf(
				"SearchUsersForMeeting failed meeting_id=%d user_id=%d keyword=%q err=%v",
				meetingID,
				userID,
				keyword,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "ユーザーの検索に失敗しました",
				},
			)
		}

		return c.JSON(http.StatusOK, users)
	}
}

// ============================================
// ユーザー削除
// ============================================
func DeactivateCurrentUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getAuthenticatedUserID(c)
		if err != nil {
			return err
		}

		err = repository.DeactivateUser(db, userID)
		if err != nil {
			switch {
			case errors.Is(
				err,
				repository.ErrUserHasIncompleteMeetings,
			):
				return c.JSON(
					http.StatusConflict,
					map[string]string{
						"error": "主催している未完了の会議があります。すべて完了してから退会してください",
					},
				)

			case errors.Is(
				err,
				repository.ErrUserNotFound,
			):
				return c.JSON(
					http.StatusNotFound,
					map[string]string{
						"error": "ユーザーが見つかりません",
					},
				)
			}

			log.Printf(
				"DeactivateCurrentUser failed user_id=%d err=%v",
				userID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "退会処理に失敗しました",
				},
			)
		}

		// 認証Cookieを削除
		cookieSecure := os.Getenv("COOKIE_SECURE") == "true"
		sameSite := http.SameSiteLaxMode
		if cookieSecure {
				sameSite = http.SameSiteNoneMode
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   cookieSecure,
			SameSite: sameSite,
			MaxAge:   -1,
		})

		return c.NoContent(http.StatusNoContent)
	}
}