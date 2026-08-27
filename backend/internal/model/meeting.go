package model

import "time"

type Meeting struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	TargetName       string     `json:"target_name"`
	ScheduledStartAt *time.Time `json:"scheduled_start_at"`
	PlannedMinutes   int        `json:"planned_minutes"`
	Decisions        string     `json:"decisions"`
	Todo             string     `json:"todo"`
	Status           string     `json:"status"`
}

type Agenda struct {
	ID               int64  `json:"id"`
	Title            string `json:"title"`
	Purpose          string `json:"purpose"`
	DiscussionPoints string `json:"discussion_points"`
	Questions        string `json:"questions"`
	Memo             string `json:"memo"`
	PlannedMinutes   int    `json:"planned_minutes"`
	SortOrder        int    `json:"sort_order"`
}

// ============================================
// 新規作成
// ============================================
type CreateMeetingRequest struct {
	Title            string                `json:"title"`
	Description      string                `json:"description"`
	TargetName       string                `json:"target_name"`
	ScheduledStartAt string                `json:"scheduled_start_at"`
	PlannedMinutes   int                   `json:"planned_minutes"`
	Decisions        string                `json:"decisions"`
	Todo             string                `json:"todo"`
	Agendas          []CreateAgendaRequest `json:"agendas"`
}

type CreateAgendaRequest struct {
	Title            string `json:"title"`
	Purpose          string `json:"purpose"`
	DiscussionPoints string `json:"discussion_points"`
	Questions        string `json:"questions"`
	Memo             string `json:"memo"`
	PlannedMinutes   int    `json:"planned_minutes"`
}

// ============================================
// 詳細
// ============================================
type MeetingDetail struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	TargetName       string     `json:"target_name"`
	ScheduledStartAt *time.Time `json:"scheduled_start_at"`
	PlannedMinutes   int        `json:"planned_minutes"`
	Decisions        string     `json:"decisions"`
	Todo             string     `json:"todo"`
	Status           string     `json:"status"`

	Agendas []Agenda `json:"agendas"`

	CurrentUserRole string `json:"current_user_role"`
}

// ============================================
// 更新
// ============================================
type UpdateMeetingRequest struct {
	Title            string                `json:"title"`
	Description      string                `json:"description"`
	TargetName       string                `json:"target_name"`
	ScheduledStartAt string                `json:"scheduled_start_at"`
	PlannedMinutes   int                   `json:"planned_minutes"`
	Decisions        string                `json:"decisions"`
	Todo             string                `json:"todo"`
	Agendas          []UpdateAgendaRequest `json:"agendas"`
}

type UpdateAgendaRequest struct {
	ID               *int64 `json:"id"`
	Title            string `json:"title"`
	Purpose          string `json:"purpose"`
	DiscussionPoints string `json:"discussion_points"`
	Questions        string `json:"questions"`
	Memo             string `json:"memo"`
	PlannedMinutes   int    `json:"planned_minutes"`
}

// ============================================
// 完了
// ============================================
type CompleteMeetingRequest struct {
	Decisions string                `json:"decisions"`
	Todo      string                `json:"todo"`
	Agendas   []UpdateAgendaRequest `json:"agendas"`
}

type CompleteAgendaRequest struct {
	ID   *int64 `json:"id"`
	Memo string `json:"memo"`
}
