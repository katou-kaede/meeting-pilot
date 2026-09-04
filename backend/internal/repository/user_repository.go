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
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUserHasIncompleteMeetings = errors.New("user has incomplete meetings")
)

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
func GetUserByEmail(db *sql.DB, email string) (*model.User, error) {
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
func GetUserByID(db *sql.DB, id int64) (*model.User, error) {
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

// ============================================
// ユーザー検索
// ============================================
func SearchUsersForMeeting(
	db *sql.DB,
	meetingID int64,
	keyword string,
) ([]model.User, error) {
	searchKeyword := strings.TrimSpace(keyword)

	rows, err := db.Query(`
		SELECT
			u.id,
			u.name,
			u.email,
			u.password_hash,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE
			u.is_active = TRUE

			-- 会議のownerを除外
			AND NOT EXISTS (
				SELECT 1
				FROM meetings m
				WHERE
					m.id = $1
					AND m.created_by = u.id
			)

			-- 追加済みメンバーを除外
			AND NOT EXISTS (
				SELECT 1
				FROM meeting_members mm
				WHERE
					mm.meeting_id = $1
					AND mm.user_id = u.id
			)

			-- 氏名またはメールアドレスで部分一致
			-- AND (
			-- 	u.name ILIKE '%' || $2 || '%'
			-- 	OR u.email ILIKE '%' || $2 || '%'
			-- )

			-- メールアドレスの完全一致
			AND LOWER(u.email) = LOWER($2)
		ORDER BY u.name
		LIMIT 20
	`,
		meetingID,
		searchKeyword,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ============================================
// ユーザー削除
// ============================================
func DeactivateUser(db *sql.DB, userID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 対象ユーザーを取得してロック
	var isActive bool

	err = tx.QueryRow(`
		SELECT is_active
		FROM users
		WHERE id = $1
		FOR UPDATE
	`,
		userID,
	).Scan(&isActive)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	if !isActive {
		return ErrUserNotFound
	}

	// 主催している未完了会議があるか確認
	var hasIncompleteMeetings bool

	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM meetings
			WHERE
				created_by = $1
				AND status <> 'completed'
		)
	`,
		userID,
	).Scan(&hasIncompleteMeetings)

	if err != nil {
		return err
	}

	if hasIncompleteMeetings {
		return ErrUserHasIncompleteMeetings
	}

	// ユーザーを論理削除
	result, err := tx.Exec(`
		UPDATE users
		SET
			is_active = FALSE,
			updated_at = NOW()
		WHERE
			id = $1
			AND is_active = TRUE
	`,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return tx.Commit()
}