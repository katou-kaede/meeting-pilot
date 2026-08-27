package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func TestLogout(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/logout",
		nil,
	)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := Logout()

	if err := handler(c); err != nil {
		t.Fatalf(
			"Logout実行時にエラーが発生しました: %v",
			err,
		)
	}

	// ステータスコードを確認
	if rec.Code != http.StatusNoContent {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusNoContent,
		)
	}

	// Set-Cookieヘッダーを確認
	cookies := rec.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal(
			"削除用Cookieが設定されていません",
		)
	}

	var accessTokenCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessTokenCookie = cookie
			break
		}
	}

	if accessTokenCookie == nil {
		t.Fatal(
			"access_token Cookieが見つかりません",
		)
	}

	if accessTokenCookie.Value != "" {
		t.Errorf(
			"Cookie Value = %q, want empty",
			accessTokenCookie.Value,
		)
	}

	if accessTokenCookie.MaxAge != -1 {
		t.Errorf(
			"Cookie MaxAge = %d, want -1",
			accessTokenCookie.MaxAge,
		)
	}

	if !accessTokenCookie.HttpOnly {
		t.Error(
			"CookieのHttpOnlyがfalseです",
		)
	}

	if accessTokenCookie.Path != "/" {
		t.Errorf(
			"Cookie Path = %q, want /",
			accessTokenCookie.Path,
		)
	}
}


const getUserByEmailSQL = `
	SELECT
		id,
		name,
		email,
		password_hash,
		is_active,
		created_at,
		updated_at
	FROM users
	WHERE LOWER(email) = $1
`

func newLoginContext(
	e *echo.Echo,
	body string,
) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(body),
	)
	req.Header.Set(
		echo.HeaderContentType,
		echo.MIMEApplicationJSON,
	)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func createPasswordHash(
	t *testing.T,
	password string,
) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf(
			"パスワードハッシュの生成に失敗しました: %v",
			err,
		)
	}

	return string(hash)
}

func TestLoginSuccess(t *testing.T) {
	t.Setenv(
		"JWT_SECRET",
		"test-secret-key-for-login-handler",
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	passwordHash := createPasswordHash(
		t,
		"password123",
	)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at",
	}).AddRow(
		int64(1),
		"Employee1",
		"employee1@example.com",
		passwordHash,
		true,
		now,
		now,
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(getUserByEmailSQL),
	).
		WithArgs("employee1@example.com").
		WillReturnRows(rows)

	e := echo.New()

	c, rec := newLoginContext(
		e,
		`{
			"email": "employee1@example.com",
			"password": "password123"
		}`,
	)

	if err := Login(db)(c); err != nil {
		t.Fatalf(
			"Login実行時にエラーが発生しました: %v",
			err,
		)
	}

	if rec.Code != http.StatusOK {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	cookies := rec.Result().Cookies()

	var accessTokenCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessTokenCookie = cookie
			break
		}
	}

	if accessTokenCookie == nil {
		t.Fatal(
			"access_token Cookieが設定されていません",
		)
	}

	if accessTokenCookie.Value == "" {
		t.Error(
			"access_token Cookieが空です",
		)
	}

	if !accessTokenCookie.HttpOnly {
		t.Error(
			"access_token CookieのHttpOnlyがfalseです",
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(getUserByEmailSQL),
	).
		WithArgs("unknown@example.com").
		WillReturnError(sql.ErrNoRows)

	e := echo.New()

	c, rec := newLoginContext(
		e,
		`{
			"email": "unknown@example.com",
			"password": "password123"
		}`,
	)

	if err := Login(db)(c); err != nil {
		t.Fatalf(
			"Login実行時にエラーが発生しました: %v",
			err,
		)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusUnauthorized,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"メールアドレスまたはパスワードが正しくありません",
	) {
		t.Errorf(
			"想定したエラーメッセージではありません: %s",
			rec.Body.String(),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestLoginPasswordMismatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	passwordHash := createPasswordHash(
		t,
		"password123",
	)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at",
	}).AddRow(
		int64(1),
		"Employee1",
		"employee1@example.com",
		passwordHash,
		true,
		now,
		now,
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(getUserByEmailSQL),
	).
		WithArgs("employee1@example.com").
		WillReturnRows(rows)

	e := echo.New()

	c, rec := newLoginContext(
		e,
		`{
			"email": "employee1@example.com",
			"password": "wrong-password"
		}`,
	)

	if err := Login(db)(c); err != nil {
		t.Fatalf(
			"Login実行時にエラーが発生しました: %v",
			err,
		)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusUnauthorized,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestLoginInactiveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	passwordHash := createPasswordHash(
		t,
		"password123",
	)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at",
	}).AddRow(
		int64(1),
		"Employee1",
		"employee1@example.com",
		passwordHash,
		false,
		now,
		now,
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(getUserByEmailSQL),
	).
		WithArgs("employee1@example.com").
		WillReturnRows(rows)

	e := echo.New()

	c, rec := newLoginContext(
		e,
		`{
			"email": "employee1@example.com",
			"password": "password123"
		}`,
	)

	if err := Login(db)(c); err != nil {
		t.Fatalf(
			"Login実行時にエラーが発生しました: %v",
			err,
		)
	}

	if rec.Code != http.StatusForbidden {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusForbidden,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"このアカウントは現在利用できません",
	) {
		t.Errorf(
			"想定したエラーメッセージではありません: %s",
			rec.Body.String(),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestLoginDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	mock.ExpectQuery(
		regexp.QuoteMeta(getUserByEmailSQL),
	).
		WithArgs("employee1@example.com").
		WillReturnError(
			errors.New("database connection failed"),
		)

	e := echo.New()

	c, rec := newLoginContext(
		e,
		`{
			"email": "employee1@example.com",
			"password": "password123"
		}`,
	)

	if err := Login(db)(c); err != nil {
		t.Fatalf(
			"Login実行時にエラーが発生しました: %v",
			err,
		)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf(
			"StatusCode = %d, want %d",
			rec.Code,
			http.StatusInternalServerError,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}