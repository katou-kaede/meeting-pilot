package model

import "time"

// 会議メンバーの取得結果
type MeetingMember struct {
	MeetingID int64     `json:"meeting_id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// 会議メンバー追加リクエスト
type CreateMeetingMemberRequest struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

// Session編集者の変更リクエスト
type UpdateMeetingEditorRequest struct {
	UserID *int64 `json:"user_id"`
}
