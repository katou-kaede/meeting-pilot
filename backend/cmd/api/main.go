package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"meeting-pilot/internal/database"
	"meeting-pilot/internal/handler"
	wsHub "meeting-pilot/internal/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	hub := wsHub.NewHub()

	// ============================================
	// APIルーティング
	// ============================================
	e.GET("/ws/meetings/:id", handler.MeetingWebSocket(hub)) // WebSocket接続

	e.GET("/api/meetings", handler.GetMeetings(db))          // 一覧
	e.GET("/api/meetings/:id", handler.GetMeetingByID(db))   // 詳細
	e.POST("/api/meetings", handler.CreateMeeting(db))       // 新規作成
	e.PUT("/api/meetings/:id", handler.UpdateMeeting(db))    // 更新
	e.DELETE("/api/meetings/:id", handler.DeleteMeeting(db)) // 削除

	e.PATCH("/api/meetings/:id/start", handler.StartMeeting(db, hub))                 // 会議開始
	e.PATCH("/api/meetings/:id/complete", handler.CompleteMeeting(db, hub))           // 会議終了
	e.GET("/api/meetings/:id/session", handler.GetMeetingSessionByID(db))             // 会議中：詳細
	e.PATCH("/api/meetings/:id/current-agenda", handler.ChangeCurrentAgenda(db, hub)) // 会議中：議題を戻す/進める
	e.PATCH("/api/meetings/:id/session", handler.SaveMeetingSession(db))              // 会議中：一時保存

	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
