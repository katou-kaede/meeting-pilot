package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	wsHub "meeting-pilot/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:5173"
		}
		// 指定したフロントエンドのURLからの接続のみ許可する
		return r.Header.Get("Origin") == allowedOrigin
	},
}

func MeetingWebSocket(db *sql.DB, hub *wsHub.Hub) echo.HandlerFunc {
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
				"MeetingWebSocket invalid id=%s err=%v",
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

		conn, err := upgrader.Upgrade(
			c.Response(),
			c.Request(),
			nil,
		)
		if err != nil {
			log.Printf(
				"MeetingWebSocket upgrade failed meeting_id=%d err=%v",
				meetingID,
				err,
			)
			// UpgraderがHTTPエラーを書き込むため、二重応答を避ける
			return nil
		}

		hub.AddClient(meetingID, conn)
		defer hub.RemoveClient(meetingID, conn)

		log.Printf(
			"MeetingWebSocket connected meeting_id=%d user_id=%d",
			meetingID,
			userID,
		)

		// 接続が切れるまで待機する
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) {
					log.Printf(
						"MeetingWebSocket read failed meeting_id=%d err=%v",
						meetingID,
						err,
					)
				}

				break
			}
		}

		log.Printf(
			"MeetingWebSocket disconnected meeting_id=%d user_id=%d",
			meetingID,
			userID,
		)

		return nil
	}
}
