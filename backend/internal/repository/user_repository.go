package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"meeting-pilot/internal/model"
)

// ============================================
// 共通エラー
// ============================================
var ErrEmailAlreadyExists = errors.New("email already exists")

// ============================================
// ユーザー登録
// ============================================
func CreateUser(
	db *sql.DB,
	req model.CreateUserRequest,
) (*model.User, error) {
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	var user model.User

	err = db.QueryRow(`
		INSERT INTO users (
			name,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			name,
			email,
			password_hash,
			is_active,
			created_at,
			updated_at
	`,
		strings.TrimSpace(req.Name),
		email,
		string(passwordHash),
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" && // 一意制約
			pgErr.ConstraintName == "uq_users_email_lower" {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

// ============================================
// ログイン
// ============================================
func GetUserByEmail(
	db *sql.DB,
	email string,
) (*model.User, error) {
	normalizedEmail := strings.ToLower(
		strings.TrimSpace(email),
	)

	var user model.User

	err := db.QueryRow(`
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
	`,
		normalizedEmail,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ============================================
// ログイン後の認証ユーザー情報取得
// ============================================
func GetUserByID(
	db *sql.DB,
	id int64,
) (*model.User, error) {
	var user model.User

	err := db.QueryRow(`
		SELECT
			id,
			name,
			email,
			password_hash,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
