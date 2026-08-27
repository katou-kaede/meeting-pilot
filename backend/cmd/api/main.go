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
	appMiddleware "meeting-pilot/internal/middleware"
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
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{
				"http://localhost:5173",
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
			},
			AllowCredentials: true,
		},
	))

	hub := wsHub.NewHub()

	// ============================================
	// APIルーティング
	// ============================================
	e.GET("/ws/meetings/:id", handler.MeetingWebSocket(db, hub), appMiddleware.RequireAuth()) // WebSocket接続

	e.GET("/api/meetings", handler.GetMeetings(db), appMiddleware.RequireAuth())          // 一覧
	e.GET("/api/meetings/:id", handler.GetMeetingByID(db), appMiddleware.RequireAuth())   // 詳細
	e.POST("/api/meetings", handler.CreateMeeting(db), appMiddleware.RequireAuth())       // 新規作成
	e.PUT("/api/meetings/:id", handler.UpdateMeeting(db), appMiddleware.RequireAuth())    // 更新
	e.DELETE("/api/meetings/:id", handler.DeleteMeeting(db), appMiddleware.RequireAuth()) // 削除

	e.PATCH("/api/meetings/:id/start",
		handler.StartMeeting(db, hub), appMiddleware.RequireAuth()) // 会議開始
	e.PATCH("/api/meetings/:id/complete",
		handler.CompleteMeeting(db, hub), appMiddleware.RequireAuth()) // 会議終了
	e.GET("/api/meetings/:id/session",
		handler.GetMeetingSessionByID(db), appMiddleware.RequireAuth()) // 会議中：詳細
	e.PATCH("/api/meetings/:id/current-agenda",
		handler.ChangeCurrentAgenda(db, hub), appMiddleware.RequireAuth()) // 会議中：議題を戻す/進める
	e.PATCH("/api/meetings/:id/session",
		handler.SaveMeetingSession(db, hub), appMiddleware.RequireAuth()) // 会議中：一時保存
	e.PATCH("/api/meetings/:id/pause",
		handler.PauseMeeting(db, hub), appMiddleware.RequireAuth()) // 会議中：一時停止
	e.PATCH("/api/meetings/:id/resume",
		handler.ResumeMeeting(db, hub), appMiddleware.RequireAuth()) // 会議中：再開

	e.POST("/api/users", handler.CreateUser(db))                              // ユーザー登録
	e.POST("/api/login", handler.Login(db))                                   // ログイン
	e.GET("/api/me", handler.GetCurrentUser(db), appMiddleware.RequireAuth()) // ログイン後の認証ユーザー情報取得
	e.POST("/api/logout", handler.Logout())                                   // ログアウト

	e.GET("/api/meetings/:id/members",
		handler.GetMeetingMembers(db), appMiddleware.RequireAuth()) // 会議参加メンバー取得
	e.GET("/api/meetings/:id/member-candidates",
		handler.SearchUsersForMeeting(db), appMiddleware.RequireAuth()) // ユーザー検索
	e.POST("/api/meetings/:id/members",
		handler.CreateMeetingMember(db), appMiddleware.RequireAuth()) // 会議メンバー追加
	e.DELETE("/api/meetings/:id/members/:userId",
		handler.DeleteMeetingMember(db), appMiddleware.RequireAuth()) // 会議メンバー削除
	e.PATCH("/api/meetings/:id/editor",
		handler.UpdateMeetingEditor(db, hub), appMiddleware.RequireAuth()) // 編集者変更

	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
