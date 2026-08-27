package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	appMiddleware "meeting-pilot/internal/middleware"
	"meeting-pilot/internal/repository"
)

// ============================================
// ログイン確認
// ============================================
func getAuthenticatedUserID(
	c echo.Context,
) (int64, error) {
	userID, ok := c.Get(
		appMiddleware.UserIDContextKey,
	).(int64)

	if !ok {
		return 0, c.JSON(
			http.StatusUnauthorized,
			map[string]string{
				"error": "ログインが必要です",
			},
		)
	}

	return userID, nil
}

// ============================================
// 会議参加メンバーかどうかチェック
// ============================================
func getMeetingUserRole(
	c echo.Context,
	db *sql.DB,
	meetingID int64,
	userID int64,
) (string, error) {
	role, err := repository.GetMeetingUserRole(
		db,
		meetingID,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "この会議を閲覧する権限がありません",
				},
			)
		}

		log.Printf(
			"meeting role check failed meeting_id=%d user_id=%d err=%v",
			meetingID,
			userID,
			err,
		)

		return "", c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "会議の権限確認に失敗しました",
			},
		)
	}

	return role, nil
}

// ============================================
// 編集可能かどうかチェック
// ============================================
func canEditMeetingSession(
	db *sql.DB,
	meetingID int64,
	userID int64,
	role string,
) (bool, error) {
	editorUserID, err :=
		repository.GetMeetingEditorUserID(
			db,
			meetingID,
		)
	if err != nil {
		return false, err
	}

	return hasMeetingSessionEditPermission(
		role,
		userID,
		editorUserID,
	), nil
}

func hasMeetingSessionEditPermission(
	role string,
	userID int64,
	editorUserID *int64,
) bool {
	// editor設定済みなら、そのeditor本人だけ編集可能
	if editorUserID != nil {
		return *editorUserID == userID
	}

	// editor未設定ならownerだけ編集可能
	return role == "owner"
}