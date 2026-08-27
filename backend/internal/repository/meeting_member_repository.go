package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"meeting-pilot/internal/model"
)

// ============================================
// 会議参加メンバー取得
// ============================================
func GetMeetingMembers(
	db *sql.DB,
	meetingID int64,
) ([]model.MeetingMember, error) {
	rows, err := db.Query(`
		SELECT
			member.meeting_id,
			member.user_id,
			member.name,
			member.email,
			member.role,
			member.created_at
		FROM (
			-- 会議作成者
			SELECT
				m.id AS meeting_id,
				u.id AS user_id,
				u.name,
				u.email,
				'owner' AS role,
				m.created_at
			FROM meetings m
			INNER JOIN users u
				ON u.id = m.created_by
			WHERE m.id = $1

			UNION ALL

			-- 追加メンバー
			SELECT
				mm.meeting_id,
				u.id AS user_id,
				u.name,
				u.email,
				mm.role,
				mm.created_at
			FROM meeting_members mm
			INNER JOIN users u
				ON u.id = mm.user_id
			WHERE mm.meeting_id = $1
		) AS member
		ORDER BY
			CASE member.role
				WHEN 'owner' THEN 1
				WHEN 'editor' THEN 2
				ELSE 3
			END,
			member.name
	`,
		meetingID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]model.MeetingMember, 0)

	for rows.Next() {
		var member model.MeetingMember

		if err := rows.Scan(
			&member.MeetingID,
			&member.UserID,
			&member.Name,
			&member.Email,
			&member.Role,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// ============================================
// 会議参加メンバーの権限取得
// ============================================
// 会議に関係のないユーザーの場合はsql.ErrNoRowsが返ります
func GetMeetingUserRole(
	db *sql.DB,
	meetingID int64,
	userID int64,
) (string, error) {
	var role string

	err := db.QueryRow(`
		SELECT role
		FROM (
			SELECT
				'owner' AS role
			FROM meetings
			WHERE
				id = $1
				AND created_by = $2

			UNION ALL

			SELECT
				mm.role
			FROM meeting_members mm
			WHERE
				mm.meeting_id = $1
				AND mm.user_id = $2
		) AS meeting_roles
		LIMIT 1
	`,
		meetingID,
		userID,
	).Scan(&role)

	if err != nil {
		return "", err
	}

	return role, nil
}

// ============================================
// 会議参加メンバーの追加
// ============================================
var (
	ErrMeetingMemberAlreadyExists = errors.New("meeting member already exists")

	ErrMeetingEditorAlreadyExists = errors.New("meeting editor already exists")

	ErrMeetingMemberCannotBeAdded = errors.New("meeting member cannot be added")
)

func CreateMeetingMember(
	db *sql.DB,
	meetingID int64,
	req model.CreateMeetingMemberRequest,
) error {
	result, err := db.Exec(`
		INSERT INTO meeting_members (
			meeting_id,
			user_id,
			role
		)
		SELECT
			$1,
			u.id,
			$2
		FROM users u
		WHERE
			u.id = $3
			AND u.is_active = TRUE

			-- ownerはmeeting_membersへ追加しない
			AND NOT EXISTS (
				SELECT 1
				FROM meetings m
				WHERE
					m.id = $1
					AND m.created_by = u.id
			)
	`,
		meetingID,
		req.Role,
		req.UserID,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" { // 一意制約

			switch pgErr.ConstraintName {
			case "meeting_members_pkey":
				return ErrMeetingMemberAlreadyExists

			case "uq_meeting_members_single_editor":
				return ErrMeetingEditorAlreadyExists
			}
		}

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrMeetingMemberCannotBeAdded
	}

	return nil
}

// ============================================
// 会議参加メンバーの削除
// ============================================
var ErrMeetingMemberNotFound = errors.New("meeting member not found")

func DeleteMeetingMember(
	db *sql.DB,
	meetingID int64,
	memberUserID int64,
) error {
	result, err := db.Exec(`
		DELETE FROM meeting_members
		WHERE
			meeting_id = $1
			AND user_id = $2
	`,
		meetingID,
		memberUserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrMeetingMemberNotFound
	}

	return nil
}

// ============================================
// 編集者の変更
// ============================================
var ErrMeetingEditorCannotBeUpdated = errors.New("meeting editor cannot be updated")

func UpdateMeetingEditor(
	db *sql.DB,
	meetingID int64,
	req model.UpdateMeetingEditorRequest,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 進行中の会議か確認
	var exists bool

	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM meetings
			WHERE
				id = $1
				AND status = 'in_progress'
		)
	`,
		meetingID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrMeetingEditorCannotBeUpdated
	}

	// 現在の編集者を参加者へ戻す
	_, err = tx.Exec(`
		UPDATE meeting_members
		SET role = 'viewer'
		WHERE
			meeting_id = $1
			AND role = 'editor'
	`,
		meetingID,
	)
	if err != nil {
		return err
	}

	// user_idがnullなら、編集者解除だけで終了
	if req.UserID == nil {
		return tx.Commit()
	}

	// 指定された参加者を編集者へ変更
	result, err := tx.Exec(`
		UPDATE meeting_members
		SET role = 'editor'
		WHERE
			meeting_id = $1
			AND user_id = $2
	`,
		meetingID,
		*req.UserID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrMeetingEditorCannotBeUpdated
	}

	return tx.Commit()
}

// ============================================
// 編集者の有無チェック
// ============================================
func GetMeetingEditorUserID(
	db *sql.DB,
	meetingID int64,
) (*int64, error) {
	var editorUserID int64

	err := db.QueryRow(`
		SELECT user_id
		FROM meeting_members
		WHERE
			meeting_id = $1
			AND role = 'editor'
	`,
		meetingID,
	).Scan(&editorUserID)

	if errors.Is(err, sql.ErrNoRows) {
		// editor未設定
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &editorUserID, nil
}
