package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"meeting-pilot/internal/model"
	"meeting-pilot/internal/repository"

	"github.com/labstack/echo/v4"
)

// ============================================
// ミーティング開始
// ============================================
func StartMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(
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

		err = repository.StartMeeting(db, id)
		if err != nil {
			log.Printf(
                "StartMeeting failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の開始に失敗しました",
				},
			)
		}

		log.Printf(
            "StartMeeting success id=%d",
            id,
        )

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// ミーティング終了
// ============================================
func CompleteMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(
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

		var req model.CompleteMeetingRequest

		if err := c.Bind(&req); err != nil {
			log.Printf(
                "CompleteMeeting bind failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "不正なリクエストです",
				},
			)
		}

		if err := repository.CompleteMeeting(db, id, req); err != nil {
			log.Printf(
                "CompleteMeeting failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の終了に失敗しました",
				},
			)
		}

		log.Printf("CompleteMeeting success id=%d", id)

		return c.NoContent(http.StatusNoContent)
	}
}

// ============================================
// 会議中：ミーティング詳細取得
// ============================================
func GetMeetingSessionByID(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(
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

		meeting, err := repository.GetMeetingSessionByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf(
                    "GetMeetingSessionByID not found id=%d",
                    id,
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
                id,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議情報の取得に失敗しました",
				},
			)
		}

		return c.JSON(http.StatusOK, meeting)
	}
}

// ============================================
// 会議中：議題を戻す/進める
// ============================================
func ChangeCurrentAgenda(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		log.Printf(
            "ChangeCurrentAgenda success meetingID=%d targetAgendaID=%d",
            meetingID,
            req.TargetAgendaID,
        )
		return c.NoContent(http.StatusNoContent)
	}
}