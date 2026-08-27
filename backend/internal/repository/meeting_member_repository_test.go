package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"

	"meeting-pilot/internal/model"
)

const getMeetingUserRoleSQL = `
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
`

func TestGetMeetingUserRole(t *testing.T) {
	tests := []struct {
		name     string
		meetingID int64
		userID    int64
		role      string
		queryErr  error
		wantRole  string
		wantErr   error
	}{
		{
			name:      "作成者ならownerを取得",
			meetingID: 1,
			userID:    10,
			role:      "owner",
			wantRole:  "owner",
		},
		{
			name:      "編集者ならeditorを取得",
			meetingID: 1,
			userID:    20,
			role:      "editor",
			wantRole:  "editor",
		},
		{
			name:      "通常参加者ならviewerを取得",
			meetingID: 1,
			userID:    30,
			role:      "viewer",
			wantRole:  "viewer",
		},
		{
			name:      "会議関係者でなければErrNoRows",
			meetingID: 1,
			userID:    40,
			queryErr:  sql.ErrNoRows,
			wantErr:   sql.ErrNoRows,
		},
		{
			name:      "DBエラー",
			meetingID: 1,
			userID:    50,
			queryErr:  errors.New("database error"),
			wantErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			expectation := mock.ExpectQuery(
				regexp.QuoteMeta(getMeetingUserRoleSQL),
			).WithArgs(
				tt.meetingID,
				tt.userID,
			)

			if tt.queryErr != nil {
				expectation.WillReturnError(tt.queryErr)
			} else {
				rows := sqlmock.NewRows(
					[]string{"role"},
				).AddRow(tt.role)

				expectation.WillReturnRows(rows)
			}

			gotRole, err := GetMeetingUserRole(
				db,
				tt.meetingID,
				tt.userID,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal(
						"エラーを期待しましたがnilでした",
					)
				}

				if !errors.Is(err, tt.wantErr) &&
					err.Error() != tt.wantErr.Error() {
					t.Errorf(
						"error = %v, want %v",
						err,
						tt.wantErr,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"エラーを期待していませんが発生しました: %v",
					err,
				)
			}

			if gotRole != tt.wantRole {
				t.Errorf(
					"role = %q, want %q",
					gotRole,
					tt.wantRole,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}

const getMeetingEditorUserIDSQL = `
	SELECT user_id
	FROM meeting_members
	WHERE
		meeting_id = $1
		AND role = 'editor'
`

func TestGetMeetingEditorUserID(t *testing.T) {
	databaseError := errors.New("database error")

	tests := []struct {
		name         string
		meetingID    int64
		editorUserID int64
		queryErr     error
		wantUserID   *int64
		wantErr      error
	}{
		{
			name:         "editorが設定されている",
			meetingID:    1,
			editorUserID: 20,
			wantUserID:   int64Pointer(20),
		},
		{
			name:      "editorが設定されていない",
			meetingID: 1,
			queryErr:  sql.ErrNoRows,
			wantUserID: nil,
		},
		{
			name:      "DBエラー",
			meetingID: 1,
			queryErr:  databaseError,
			wantErr:   databaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			expectation := mock.ExpectQuery(
				regexp.QuoteMeta(
					getMeetingEditorUserIDSQL,
				),
			).WithArgs(tt.meetingID)

			if tt.queryErr != nil {
				expectation.WillReturnError(tt.queryErr)
			} else {
				expectation.WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id"},
					).AddRow(tt.editorUserID),
				)
			}

			got, err := GetMeetingEditorUserID(
				db,
				tt.meetingID,
			)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"error = %v, want %v",
						err,
						tt.wantErr,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"エラーを期待していませんが発生しました: %v",
					err,
				)
			}

			if tt.wantUserID == nil {
				if got != nil {
					t.Errorf(
						"editorUserID = %v, want nil",
						*got,
					)
				}

				return
			}

			if got == nil {
				t.Fatal(
					"editorUserIDを期待しましたがnilでした",
				)
			}

			if *got != *tt.wantUserID {
				t.Errorf(
					"editorUserID = %d, want %d",
					*got,
					*tt.wantUserID,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}

func TestCreateMeetingMemberDuplicateErrors(
	t *testing.T,
) {
	tests := []struct {
		name           string
		constraintName string
		wantErr        error
	}{
		{
			name:           "同じメンバーを重複追加するとエラー",
			constraintName: "meeting_members_pkey",
			wantErr: ErrMeetingMemberAlreadyExists,
		},
		{
			name:           "editorを2人設定するとエラー",
			constraintName: "uq_meeting_members_single_editor",
			wantErr: ErrMeetingEditorAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			req := model.CreateMeetingMemberRequest{
				UserID: 20,
				Role:   "editor",
			}

			mock.ExpectExec(
				`INSERT INTO meeting_members`,
			).
				WithArgs(
					int64(1),
					req.Role,
					req.UserID,
				).
				WillReturnError(
					&pgconn.PgError{
						Code:           "23505",
						ConstraintName: tt.constraintName,
					},
				)

			err = CreateMeetingMember(
				db,
				1,
				req,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}

func TestDeleteMeetingMember(t *testing.T) {
	databaseError := errors.New("database error")

	tests := []struct {
		name         string
		meetingID    int64
		memberUserID int64
		result       sql.Result
		execErr      error
		wantErr      error
	}{
		{
			name:         "メンバーを正常に削除できる",
			meetingID:    1,
			memberUserID: 20,
			result:       sqlmock.NewResult(0, 1),
			wantErr:      nil,
		},
		{
			name:         "対象メンバーが存在しない",
			meetingID:    1,
			memberUserID: 99,
			result:       sqlmock.NewResult(0, 0),
			wantErr:      ErrMeetingMemberNotFound,
		},
		{
			name:         "DBエラー",
			meetingID:    1,
			memberUserID: 20,
			execErr:      databaseError,
			wantErr:      databaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			expectation := mock.ExpectExec(
				`DELETE FROM meeting_members`,
			).WithArgs(
				tt.meetingID,
				tt.memberUserID,
			)

			if tt.execErr != nil {
				expectation.WillReturnError(tt.execErr)
			} else {
				expectation.WillReturnResult(tt.result)
			}

			err = DeleteMeetingMember(
				db,
				tt.meetingID,
				tt.memberUserID,
			)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf(
						"エラーを期待していませんが発生しました: %v",
						err,
					)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}