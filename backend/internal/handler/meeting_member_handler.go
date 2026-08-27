package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"meeting-pilot/internal/model"
	"meeting-pilot/internal/repository"
	wsHub "meeting-pilot/internal/websocket"
)

// ============================================
// 会議参加メンバー取得
// ============================================
func GetMeetingMembers(
	db *sql.DB,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ログインユーザーIDを取得
		userID, err := getAuthenticatedUserID(c)
		if err != nil {
			return err
		}

		// URLから会議IDを取得
		meetingID, err := strconv.ParseInt(
			c.Param("id"),
			10,
			64,
		)
		if err != nil {
			log.Printf(
				"GetMeetingMembers invalid meeting_id=%s user_id=%d err=%v",
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

		members, err := repository.GetMeetingMembers(db, meetingID)
		if err != nil {
			log.Printf(
				"GetMeetingMembers failed meeting_id=%d user_id=%d err=%v",
				meetingID,
				userID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議メンバーの取得に失敗しました",
				},
			)
		}

		return c.JSON(http.StatusOK, members)
	}
}

// ============================================
// 会議参加メンバーの追加
// ============================================
func CreateMeetingMember(
	db *sql.DB,
) echo.HandlerFunc {
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
				"CreateMeetingMember invalid meeting_id=%s user_id=%d err=%v",
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

		var req model.CreateMeetingMemberRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"CreateMeetingMember bind failed meeting_id=%d err=%v",
				meetingID,
				err,
			)

			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "入力内容の形式が正しくありません",
				},
			)
		}

		if req.UserID <= 0 {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "ユーザーIDが不正です",
				},
			)
		}

		if req.Role != "editor" &&
			req.Role != "viewer" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "ロールが不正です",
				},
			)
		}

		err = repository.CreateMeetingMember(
			db,
			meetingID,
			req,
		)
		if err != nil {
			switch {
			case errors.Is(
				err,
				repository.ErrMeetingMemberAlreadyExists,
			):
				return c.JSON(
					http.StatusConflict,
					map[string]string{
						"error": "このユーザーは既に追加されています",
					},
				)

			case errors.Is(
				err,
				repository.ErrMeetingEditorAlreadyExists,
			):
				return c.JSON(
					http.StatusConflict,
					map[string]string{
						"error": "編集者は既に設定されています",
					},
				)

			case errors.Is(
				err,
				repository.ErrMeetingMemberCannotBeAdded,
			):
				return c.JSON(
					http.StatusBadRequest,
					map[string]string{
						"error": "指定したユーザーを追加できません",
					},
				)
			}

			log.Printf(
				"CreateMeetingMember failed meeting_id=%d target_user_id=%d role=%s err=%v",
				meetingID,
				req.UserID,
				req.Role,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議メンバーの追加に失敗しました",
				},
			)
		}

		return c.NoContent(http.StatusCreated)
	}
}

// ============================================
// 会議参加メンバーの削除
// ============================================
func DeleteMeetingMember(
	db *sql.DB,
) echo.HandlerFunc {
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
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "会議IDが不正です",
				},
			)
		}

		memberUserID, err := strconv.ParseInt(
			c.Param("userId"),
			10,
			64,
		)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "ユーザーIDが不正です",
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

		err = repository.DeleteMeetingMember(
			db,
			meetingID,
			memberUserID,
		)
		if err != nil {
			if errors.Is(
				err,
				repository.ErrMeetingMemberNotFound,
			) {
				return c.JSON(
					http.StatusNotFound,
					map[string]string{
						"error": "指定したメンバーが見つかりません",
					},
				)
			}

			log.Printf(
				"DeleteMeetingMember failed meeting_id=%d member_user_id=%d operator_user_id=%d err=%v",
				meetingID,
				memberUserID,
				userID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議メンバーの削除に失敗しました",
				},
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 編集者の変更
// ============================================
func UpdateMeetingEditor(
	db *sql.DB,
	hub *wsHub.Hub,
) echo.HandlerFunc {
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
				"UpdateMeetingEditor invalid meeting_id=%s user_id=%d err=%v",
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

		// ownerか確認
		// ログインユーザーの権限チェック
		role, err := getMeetingUserRole(
			c,
			db,
			meetingID,
			userID,
		)
		if err != nil {
			return err
		}

		if role != "owner" {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "編集者を変更できるのは会議作成者のみです",
				},
			)
		}

		var req model.UpdateMeetingEditorRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"UpdateMeetingEditor bind failed meeting_id=%d user_id=%d err=%v",
				meetingID,
				userID,
				err,
			)

			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "入力内容の形式が正しくありません",
				},
			)
		}

		err = repository.UpdateMeetingEditor(
			db,
			meetingID,
			req,
		)
		if err != nil {
			if errors.Is(
				err,
				repository.ErrMeetingEditorCannotBeUpdated,
			) {
				return c.JSON(
					http.StatusBadRequest,
					map[string]string{
						"error": "指定した編集者を設定できません",
					},
				)
			}

			log.Printf(
				"UpdateMeetingEditor failed meeting_id=%d user_id=%d target_user_id=%v err=%v",
				meetingID,
				userID,
				req.UserID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "編集者の変更に失敗しました",
				},
			)
		}

		// Sessionを開いている他の画面へ通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_editor_changed",
				MeetingID: meetingID,
			},
		)

		return c.NoContent(http.StatusNoContent)
	}
}
