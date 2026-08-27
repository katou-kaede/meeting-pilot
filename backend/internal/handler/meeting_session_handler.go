package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"meeting-pilot/internal/model"
	"meeting-pilot/internal/repository"
	wsHub "meeting-pilot/internal/websocket"

	"github.com/labstack/echo/v4"
)

// ============================================
// ミーティング開始
// ============================================
func StartMeeting(db *sql.DB, hub *wsHub.Hub) echo.HandlerFunc {
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

		err = repository.StartMeeting(db, meetingID)
		if err != nil {
			log.Printf(
				"StartMeeting failed id=%d err=%v",
				meetingID,
				err,
			)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の開始に失敗しました",
				},
			)
		}

		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_started",
				MeetingID: meetingID,
			},
		)

		log.Printf(
			"StartMeeting success id=%d",
			meetingID,
		)

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// ミーティング終了
// ============================================
func CompleteMeeting(db *sql.DB, hub *wsHub.Hub) echo.HandlerFunc {
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
				"CompleteMeeting invalid id=%s err=%v",
				c.Param("id"),
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

		var req model.CompleteMeetingRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"CompleteMeeting bind failed id=%d err=%v",
				meetingID,
				err,
			)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "不正なリクエストです",
				},
			)
		}

		if err := repository.CompleteMeeting(db, meetingID, req); err != nil {
			log.Printf(
				"CompleteMeeting failed id=%d err=%v",
				meetingID,
				err,
			)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の終了に失敗しました",
				},
			)
		}

		// DB更新成功後に、同じ会議へ接続中のブラウザへ通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_completed",
				MeetingID: meetingID,
			},
		)

		log.Printf("CompleteMeeting success id=%d", meetingID)

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 会議中：ミーティング詳細取得
// ============================================
func GetMeetingSessionByID(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				"GetMeetingSessionByID invalid id=%s err=%v",
				c.Param("id"),
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
		role, err := getMeetingUserRole(
			c,
			db,
			meetingID,
			userID,
		)
		if err != nil {
			return err
		}

		meeting, err := repository.GetMeetingSessionByID(db, meetingID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf(
					"GetMeetingSessionByID not found id=%d",
					meetingID,
				)
				return c.JSON(
					http.StatusNotFound,
					map[string]string{
						"error": "会議が見つかりません",
					},
				)
			}

			log.Printf(
				"GetMeetingSessionByID failed id=%d err=%v",
				meetingID,
				err,
			)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議情報の取得に失敗しました",
				},
			)
		}

		meeting.CurrentUserRole = role

		// editorがいればeditorのみ編集可能
		// editorがいなければownerが編集可能
		meeting.CanEditSession =
			role == "editor" ||
				(role == "owner" && meeting.EditorUserID == nil)

		return c.JSON(http.StatusOK, meeting)
	}
}

// ============================================
// 会議中：議題を戻す/進める
// ============================================
func ChangeCurrentAgenda(db *sql.DB, hub *wsHub.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				"ChangeCurrentAgenda invalid meeting id=%s err=%v",
				c.Param("id"),
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

		var req model.ChangeCurrentAgendaRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"ChangeCurrentAgenda bind failed meetingID=%d err=%v",
				meetingID,
				err,
			)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "不正なリクエストです",
				},
			)
		}

		if req.TargetAgendaID <= 0 {
			log.Printf(
				"ChangeCurrentAgenda invalid agenda id meetingID=%d targetAgendaID=%d",
				meetingID,
				req.TargetAgendaID,
			)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "不正なアジェンダです",
				},
			)
		}

		if err := repository.ChangeCurrentAgenda(
			db,
			meetingID,
			req,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf(
					"ChangeCurrentAgenda not found meetingID=%d targetAgendaID=%d",
					meetingID,
					req.TargetAgendaID,
				)
				return c.JSON(
					http.StatusNotFound,
					map[string]string{
						"error": "会議またはアジェンダが見つかりません",
					},
				)
			}

			log.Printf(
				"ChangeCurrentAgenda failed meetingID=%d targetAgendaID=%d err=%v",
				meetingID,
				req.TargetAgendaID,
				err,
			)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "議題の切り替えに失敗しました",
				},
			)
		}

		// DB更新成功後に、同じ会議へ接続中のブラウザへ通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "current_agenda_changed",
				MeetingID: meetingID,
			},
		)

		log.Printf(
			"ChangeCurrentAgenda success meetingID=%d targetAgendaID=%d",
			meetingID,
			req.TargetAgendaID,
		)
		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 会議中：一時保存
// ============================================
func SaveMeetingSession(
	db *sql.DB,
	hub *wsHub.Hub,
) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				"SaveMeetingSession invalid id=%s err=%v",
				c.Param("id"),
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
		role, err := getMeetingUserRole(
			c,
			db,
			meetingID,
			userID,
		)
		if err != nil {
			return err
		}

		canEdit, err := canEditMeetingSession(
			db,
			meetingID,
			userID,
			role,
		)
		if err != nil {
			log.Printf(
				"SaveMeetingSession permission check failed meeting_id=%d user_id=%d err=%v",
				meetingID,
				userID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "編集権限の確認に失敗しました",
				},
			)
		}
		if !canEdit {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "セッション内容を編集する権限がありません",
				},
			)
		}

		var req model.SaveMeetingSessionRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
				"SaveMeetingSession bind failed id=%d err=%v",
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

		if err := repository.SaveMeetingSession(db, meetingID, req); err != nil {
			log.Printf(
				"SaveMeetingSession failed id=%d err=%v",
				meetingID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議内容の一時保存に失敗しました",
				},
			)
		}

		// 同じ会議を開いている画面へ通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_session_saved",
				MeetingID: meetingID,
			},
		)

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 会議中：一時停止
// ============================================
func PauseMeeting(
	db *sql.DB,
	hub *wsHub.Hub,
) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				"PauseMeeting invalid id=%s err=%v",
				c.Param("id"),
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

		if err := repository.PauseMeeting(db, meetingID); err != nil {
			log.Printf(
				"PauseMeeting failed id=%d err=%v",
				meetingID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の一時停止に失敗しました",
				},
			)
		}

		// 同じ会議を開いている画面へ通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_paused",
				MeetingID: meetingID,
			},
		)

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 会議中：再開
// ============================================
func ResumeMeeting(
	db *sql.DB,
	hub *wsHub.Hub,
) echo.HandlerFunc {
	return func(c echo.Context) error {
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
				"ResumeMeeting invalid id=%s err=%v",
				c.Param("id"),
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

		if err := repository.ResumeMeeting(db, meetingID); err != nil {
			log.Printf(
				"ResumeMeeting failed id=%d err=%v",
				meetingID,
				err,
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の再開に失敗しました",
				},
			)
		}

		// 同じ会議を開いている画面へ再開を通知
		hub.Broadcast(
			meetingID,
			wsHub.Event{
				Type:      "meeting_resumed",
				MeetingID: meetingID,
			},
		)

		return c.NoContent(http.StatusNoContent)
	}
}
