package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	wsHub "meeting-pilot/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 開発中は許可。本番では接続元を制限する
		return true
	},
}

func MeetingWebSocket(hub *wsHub.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
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
			"MeetingWebSocket disconnected meeting_id=%d",
			meetingID,
		)

		return nil
	}
}
