package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"meeting-pilot/internal/model"
	"meeting-pilot/internal/repository"
	"meeting-pilot/internal/util"
	"meeting-pilot/internal/validator"

	"github.com/labstack/echo/v4"
)

// ============================================
// ミーティング一覧取得
// ============================================
func GetMeetings(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		meetings, err := repository.GetMeetings(db)
		if err != nil {
			log.Printf("GetMeetings failed: %v", err)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議一覧の取得に失敗しました",
				},
			)
		}

		return c.JSON(http.StatusOK, meetings)
	}
}

// ============================================
// ミーティング新規作成
// ============================================
func CreateMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		var req model.CreateMeetingRequest

		// Bind：JSONを自動でGo structに変換
		if err := c.Bind(&req); err != nil {
			log.Printf("CreateMeeting bind failed err=%v", err)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "不正なリクエストです",
				},
			)
		}

		// バリデーション
		if err := validator.ValidateCreateMeeting(req); err != nil {
			log.Printf("CreateMeeting validation failed err=%v", err)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		// 時間の形式がGoとフロントで異なるのでここでパースする
		scheduledStartAt, err := util.ParseDateTimeLocal(
			req.ScheduledStartAt,
		)
		if err != nil {
			log.Printf("CreateMeeting parse datetime failed err=%v", err)
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "日時の形式が不正です",
				},
			)
		}

		if err := repository.CreateMeeting(
			db, req, scheduledStartAt,
		); err != nil {
			log.Printf(
                "CreateMeeting failed title=%s err=%v",
                req.Title,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の登録に失敗しました",
				},
			)
		}

		log.Printf(
            "CreateMeeting success title=%s agendaCount=%d",
            req.Title,
            len(req.Agendas),
        )

		return c.NoContent(http.StatusCreated)
	}
}

// ============================================
// ミーティング詳細取得
// ============================================
func GetMeetingByID(
	db *sql.DB,
) echo.HandlerFunc {
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
					"error": "invalid id",
				},
			)
		}

		meeting, err := repository.GetMeetingByID(db, id)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		return c.JSON(http.StatusOK, meeting)
	}
}

// ============================================
// ミーティング更新
// ============================================
func UpdateMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(
			c.Param("id"),
			10,
			64,
		)
		if err != nil {
			log.Printf(
                "UpdateMeeting invalid id=%s err=%v",
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

		var req model.UpdateMeetingRequest

		// Bind：JSONを自動でGo structに変換
		if err := c.Bind(&req); err != nil {
			log.Printf(
                "UpdateMeeting bind failed id=%d err=%v",
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

		// バリデーション
		if err := validator.ValidateUpdateMeeting(req); err != nil {
			log.Printf(
                "UpdateMeeting validation failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		// 時間の形式がGoとフロントで異なるのでここでパースする
		scheduledStartAt, err := util.ParseDateTimeLocal(
			req.ScheduledStartAt,
		)
		if err != nil {
			log.Printf(
                "UpdateMeeting parse datetime failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "日時の形式が不正です",
				},
			)
		}

		err = repository.UpdateMeeting(
			db,
			id,
			req,
			scheduledStartAt,
		)
		if err != nil {
			log.Printf(
                "UpdateMeeting failed id=%d title=%s err=%v",
                id,
                req.Title,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の更新に失敗しました",
				},
			)
		}

		log.Printf(
			"UpdateMeeting success id=%d title=%s",
			id,
			req.Title,
		)

		return c.NoContent(http.StatusOK)
	}
}

// ============================================
// ミーティング削除
// ============================================
func DeleteMeeting(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(
			c.Param("id"),
			10,
			64,
		)

		if err != nil {
			log.Printf(
                "DeleteMeeting invalid id=%s err=%v",
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

		err = repository.DeleteMeeting(db, id)
		if err != nil {
			log.Printf(
                "DeleteMeeting failed id=%d err=%v",
                id,
                err,
            )
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": "会議の削除に失敗しました",
				},
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

