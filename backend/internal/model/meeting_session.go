package model

import "time"

// ============================================
// 会議中：詳細
// ============================================
type MeetingSessionDetail struct {
	ID                 int64      `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	PlannedMinutes     int        `json:"planned_minutes"`
	Status             string     `json:"status"`
	ActualStartAt      *time.Time `json:"actual_start_at"`
	PausedAt           *time.Time `json:"paused_at"`
	TotalPausedSeconds int        `json:"total_paused_seconds"`
	CurrentAgendaID    *int64     `json:"current_agenda_id"`
	Decisions          string     `json:"decisions"`
	Todo               string     `json:"todo"`

	Agendas []AgendaSessionDetail `json:"agendas"`

	EditorUserID    *int64 `json:"editor_user_id"`
	CurrentUserRole string `json:"current_user_role"`
	CanEditSession  bool   `json:"can_edit_session"`
}

type AgendaSessionDetail struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	Purpose          string     `json:"purpose"`
	DiscussionPoints string     `json:"discussion_points"`
	Questions        string     `json:"questions"`
	Memo             string     `json:"memo"`
	PlannedMinutes   int        `json:"planned_minutes"`
	SortOrder        int        `json:"sort_order"`
	ActualStartAt    *time.Time `json:"actual_start_at"`
	ActualEndAt      *time.Time `json:"actual_end_at"`
	ElapsedSeconds   *int       `json:"elapsed_seconds"`
}

// ============================================
// 会議中：議題を戻す/進める
// ============================================
type ChangeCurrentAgendaRequest struct {
	TargetAgendaID int64 `json:"agenda_id"`
}

// ============================================
// 会議中：一時保存
// ============================================
type SaveAgendaSessionRequest struct {
	ID   int64  `json:"id"`
	Memo string `json:"memo"`
}

type SaveMeetingSessionRequest struct {
	Decisions string                     `json:"decisions"`
	Todo      string                     `json:"todo"`
	Agendas   []SaveAgendaSessionRequest `json:"agendas"`
}
