package repository

import (
	"database/sql"
	"time"

	"meeting-pilot/internal/model"
)

// ============================================
// ミーティング一覧取得
// ============================================
func GetMeetings(db *sql.DB) ([]model.Meeting, error) {
	rows, err := db.Query(`
		SELECT
			id,
			title,
			COALESCE(target_name, ''),
			scheduled_start_at,
			planned_minutes,
			COALESCE(decisions, ''),
			COALESCE(todo, ''),
			status
		FROM meetings
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []model.Meeting

	for rows.Next() {
		var meeting model.Meeting

		err := rows.Scan(
			&meeting.ID,
			&meeting.Title,
			&meeting.TargetName,
			&meeting.ScheduledStartAt,
			&meeting.PlannedMinutes,
			&meeting.Decisions,
			&meeting.Todo,
			&meeting.Status,
		)
		if err != nil {
			return nil, err
		}

		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

// ============================================
// ミーティング新規作成
// ============================================
func CreateMeeting(
	db *sql.DB,
	req model.CreateMeetingRequest,
	scheduledStartAt *time.Time,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// meetingを作成
	var meetingID int64

	err = tx.QueryRow(`
		INSERT INTO meetings (
			title,
			description,
			target_name,
			scheduled_start_at,
			planned_minutes,
			decisions,
			todo,
			status,
			created_by
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			'scheduled',
			1
		)
		RETURNING id
	`,
		req.Title,
		req.Description,
		req.TargetName,
		scheduledStartAt,
		req.PlannedMinutes,
		req.Decisions,
		req.Todo,
	).Scan(&meetingID)
	if err != nil {
		return err
	}

	// 子要素のAgendasを作成
	for index, agenda := range req.Agendas {

		_, err := tx.Exec(`
			INSERT INTO agendas (
				meeting_id,
				title,
				purpose,
				discussion_points,
				questions,
				memo,
				planned_minutes,
				sort_order
			)
			VALUES (
				$1,$2,$3,$4,$5,
				$6,$7,$8
			)
		`,
			meetingID,
			agenda.Title,
			agenda.Purpose,
			agenda.DiscussionPoints,
			agenda.Questions,
			agenda.Memo,
			agenda.PlannedMinutes,
			index+1,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ============================================
// ミーティング詳細取得
// ============================================
func GetMeetingByID(db *sql.DB, id int64) (*model.MeetingDetail, error) {
	var meeting model.MeetingDetail

	// meething取得
	err := db.QueryRow(`
		SELECT
			id,
			title,
			COALESCE(description, ''),
			COALESCE(target_name, ''),
			scheduled_start_at,
			planned_minutes,
			COALESCE(decisions, ''),
			COALESCE(todo, ''),
			status
		FROM meetings
		WHERE id = $1
	`,
		id,
	).Scan(
		&meeting.ID,
		&meeting.Title,
		&meeting.Description,
		&meeting.TargetName,
		&meeting.ScheduledStartAt,
		&meeting.PlannedMinutes,
		&meeting.Decisions,
		&meeting.Todo,
		&meeting.Status,
	)
	if err != nil {
		return nil, err
	}

	// agendas取得
	rows, err := db.Query(`
		SELECT
			id,
			title,
			COALESCE(purpose, ''),
			COALESCE(discussion_points, ''),
			COALESCE(questions, ''),
			COALESCE(memo, ''),
			planned_minutes,
			sort_order
		FROM agendas
		WHERE meeting_id = $1
		ORDER BY sort_order
	`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var agenda model.Agenda

		err := rows.Scan(
			&agenda.ID,
			&agenda.Title,
			&agenda.Purpose,
			&agenda.DiscussionPoints,
			&agenda.Questions,
			&agenda.Memo,
			&agenda.PlannedMinutes,
			&agenda.SortOrder,
		)
		if err != nil {
			return nil, err
		}

		meeting.Agendas = append(
			meeting.Agendas,
			agenda,
		)
	}

	return &meeting, nil
}

// ============================================
// ミーティング更新
// ============================================
func UpdateMeeting(
	db *sql.DB,
	id int64,
	req model.UpdateMeetingRequest,
	scheduledStartAt *time.Time,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// meetingを更新
	_, err = tx.Exec(`
		UPDATE meetings
		SET
			title = $1,
			description = $2,
			target_name = $3,
			scheduled_start_at = $4,
			planned_minutes = $5,
			decisions = $6,
			todo = $7
		WHERE id = $8
	`,
		req.Title,
		req.Description,
		req.TargetName,
		scheduledStartAt,
		req.PlannedMinutes,
		req.Decisions,
		req.Todo,
		id,
	)
	if err != nil {
		return err
	}

	// ----------------------------
	// 渡されてきたAgendaIDを取得
	requestAgendaIDs := make(map[int64]bool)
	for _, agenda := range req.Agendas {
		if agenda.ID != nil {
			requestAgendaIDs[*agenda.ID] = true
		}
	}

	// DBにある子要素のAgendaIDを取得
	rows, err := tx.Query(`
		SELECT id
		FROM agendas
		WHERE meeting_id = $1
	`,
		id,
	)
	if err != nil {
		return err
	}
	// defer rows.Close()

	var dbAgendaIDs []int64

	for rows.Next() {
		var agendaID int64

		err := rows.Scan(&agendaID)
		if err != nil {
			rows.Close()
			return err
		}

		dbAgendaIDs = append(dbAgendaIDs, agendaID)
	}
	rows.Close()

	// ------- DBにあるがリクエストにない子要素のAgendasを削除 -------
	for _, agendaID := range dbAgendaIDs {
		if !requestAgendaIDs[agendaID] {
			_, err := tx.Exec(`
				DELETE
				FROM agendas
				WHERE id = $1
			`,
				agendaID,
			)
			if err != nil {
				return err
			}
		}
	}

	// ------- sort_orderを一時的に負数へ退避 -------
	for _, agenda := range req.Agendas {
		if agenda.ID == nil {
			continue
		}

		_, err := tx.Exec(`
			UPDATE agendas
			SET sort_order = -id
			WHERE id = $1
		`,
			*agenda.ID,
		)
		if err != nil {
			return err
		}
	}

	// ------- 子要素のAgendasを更新 -------
	for index, agenda := range req.Agendas {
		if agenda.ID == nil {
			continue
		}

		_, err := tx.Exec(`
			UPDATE agendas
			SET
				title = $1,
				purpose = $2,
				discussion_points = $3,
				questions = $4,
				memo = $5,
				planned_minutes = $6,
				sort_order = $7,
				updated_at = NOW()
			WHERE id = $8
		`,
			agenda.Title,
			agenda.Purpose,
			agenda.DiscussionPoints,
			agenda.Questions,
			agenda.Memo,
			agenda.PlannedMinutes,
			index+1,
			*agenda.ID,
		)
		if err != nil {
			return err
		}
	}

	// ------- 子要素のAgendasを作成 -------
	for index, agenda := range req.Agendas {
		if agenda.ID != nil {
			continue
		}

		_, err := tx.Exec(`
			INSERT INTO agendas (
				meeting_id,
				title,
				purpose,
				discussion_points,
				questions,
				memo,
				planned_minutes,
				sort_order
			)
			VALUES (
				$1,$2,$3,$4,$5,
				$6,$7,$8
			)
		`,
			id,
			agenda.Title,
			agenda.Purpose,
			agenda.DiscussionPoints,
			agenda.Questions,
			agenda.Memo,
			agenda.PlannedMinutes,
			index+1,
		)

		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ============================================
// ミーティング削除
// ============================================
func DeleteMeeting(db *sql.DB, id int64) error {

	_, err := db.Exec(`
		DELETE
		FROM meetings
		WHERE id = $1
	`,
		id,
	)

	return err
}
