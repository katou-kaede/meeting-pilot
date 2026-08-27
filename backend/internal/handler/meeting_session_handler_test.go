package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"

	appMiddleware "meeting-pilot/internal/middleware"
	wsHub "meeting-pilot/internal/websocket"
)

const getMeetingUserRoleQuery = `
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

const getMeetingEditorUserIDQuery = `
	SELECT user_id
	FROM meeting_members
	WHERE
		meeting_id = $1
		AND role = 'editor'
`

func newSaveMeetingSessionContext(
	e *echo.Echo,
	userID int64,
	body string,
) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/meetings/1/session",
		strings.NewReader(body),
	)

	req.Header.Set(
		echo.HeaderContentType,
		echo.MIMEApplicationJSON,
	)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/api/meetings/:id/session")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// RequireAuthが設定した状態を再現
	c.Set(
		appMiddleware.UserIDContextKey,
		userID,
	)

	return c, rec
}

func TestSaveMeetingSessionPermission(t *testing.T) {
	tests := []struct {
		name               string
		userID             int64
		role               string
		editorUserID       *int64
		wantStatusCode     int
		expectSaveExecuted bool
	}{
		{
			name:               "editor本人なら一時保存できる",
			userID:             20,
			role:               "editor",
			editorUserID:       int64Pointer(20),
			wantStatusCode:     http.StatusNoContent,
			expectSaveExecuted: true,
		},
		{
			name:               "editor設定中のownerは一時保存できない",
			userID:             10,
			role:               "owner",
			editorUserID:       int64Pointer(20),
			wantStatusCode:     http.StatusForbidden,
			expectSaveExecuted: false,
		},
		{
			name:               "editor未設定でもviewerは一時保存できない",
			userID:             30,
			role:               "viewer",
			editorUserID:       nil,
			wantStatusCode:     http.StatusForbidden,
			expectSaveExecuted: false,
		},
		{
			name:               "editor未設定ならownerは一時保存できる",
			userID:             10,
			role:               "owner",
			editorUserID:       nil,
			wantStatusCode:     http.StatusNoContent,
			expectSaveExecuted: true,
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

			const meetingID int64 = 1

			// 会議内でのロール取得
			mock.ExpectQuery(
				regexp.QuoteMeta(
					getMeetingUserRoleQuery,
				),
			).
				WithArgs(
					meetingID,
					tt.userID,
				).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"role"},
					).AddRow(tt.role),
				)

			// 現在のeditor取得
			editorExpectation := mock.ExpectQuery(
				regexp.QuoteMeta(
					getMeetingEditorUserIDQuery,
				),
			).
				WithArgs(meetingID)

			if tt.editorUserID == nil {
				editorExpectation.WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id"},
					),
				)
			} else {
				editorExpectation.WillReturnRows(
					sqlmock.NewRows(
						[]string{"user_id"},
					).AddRow(*tt.editorUserID),
				)
			}

			if tt.expectSaveExecuted {
				mock.ExpectBegin()

				mock.ExpectExec(
					`UPDATE meetings\s+SET\s+decisions = \$1`,
				).
					WithArgs(
						"新機能を実装する",
						"仕様を確認する",
						meetingID,
					).
					WillReturnResult(
						sqlmock.NewResult(0, 1),
					)

				// 今回はagendasを空にしているため、
				// Agenda UPDATEは実行されない
				mock.ExpectCommit()
			}

			e := echo.New()
			hub := wsHub.NewHub()

			c, rec := newSaveMeetingSessionContext(
				e,
				tt.userID,
				`{
					"decisions": "新機能を実装する",
					"todo": "仕様を確認する",
					"agendas": []
				}`,
			)

			err = SaveMeetingSession(
				db,
				hub,
			)(c)

			if err != nil {
				t.Fatalf(
					"Handler実行時にエラーが発生しました: %v",
					err,
				)
			}

			if rec.Code != tt.wantStatusCode {
				t.Errorf(
					"StatusCode = %d, want %d\nBody: %s",
					rec.Code,
					tt.wantStatusCode,
					rec.Body.String(),
				)
			}

			if tt.wantStatusCode ==
				http.StatusForbidden &&
				!strings.Contains(
					rec.Body.String(),
					"セッション内容を編集する権限がありません",
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
		})
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}