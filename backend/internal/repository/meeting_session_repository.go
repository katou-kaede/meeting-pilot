package repository

import (
	"database/sql"
	"fmt"
	"time"

	"meeting-pilot/internal/model"
)

// ============================================
// ミーティング開始
// ============================================
func StartMeeting(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 先頭の議題を取得
	var firstAgendaID int64

	err = tx.QueryRow(`
		SELECT id
		FROM agendas
		WHERE meeting_id = $1
		ORDER BY sort_order
		LIMIT 1
	`, id).Scan(&firstAgendaID)

	if err != nil {
		return err
	}

	// 会議を開始
	result, err := tx.Exec(`
		UPDATE meetings
		SET
			status = 'in_progress',
			actual_start_at = NOW(),
			current_agenda_id = $1,
			updated_at = NOW()
		WHERE
			id = $2
			AND status = 'scheduled'
	`,
		firstAgendaID,
		id,
	)
	if err != nil {
		return err
	}

	// RowsAffected: 更新された行数を取得
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("meeting not found or cannot be started")
	}

	// 先頭の議題を開始
	_, err = tx.Exec(`
		UPDATE agendas
		SET
			actual_start_at = NOW(),
			actual_end_at = NULL,
			elapsed_seconds = 0,
			updated_at = NOW()
		WHERE
			id = $1
			AND meeting_id = $2
	`,
		firstAgendaID,
		id,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ============================================
// ミーティング終了
// ============================================
func CompleteMeeting(
	db *sql.DB,
	id int64,
	req model.CompleteMeetingRequest,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 会議結果・ステータスを更新
	result, err := tx.Exec(`
		UPDATE meetings
		SET
			status = $1,
			decisions = $2,
			todo = $3,
			actual_end_at = NOW(),
			updated_at = NOW()
		WHERE 
			id = $4
			AND status = 'in_progress'
	`,
		"completed",
		req.Decisions,
		req.Todo,
		id,
	)
	if err != nil {
		return err
	}

	// RowsAffected: 更新された行数を取得
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("meeting not found or cannot be completed")
	}

	// 各アジェンダのメモを更新
	for _, agenda := range req.Agendas {
		_, err = tx.Exec(`
			UPDATE agendas
			SET
				memo = $1,
				updated_at = NOW()
			WHERE
				id = $2
				AND meeting_id = $3
		`,
			agenda.Memo,
			agenda.ID,
			id,
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
// 会議中：ミーティング詳細取得
// ============================================
func GetMeetingSessionByID(
	db *sql.DB,
	id int64,
) (*model.MeetingSessionDetail, error) {
	var meeting model.MeetingSessionDetail

	err := db.QueryRow(`
		SELECT
			id,
			title,
			COALESCE(description, ''),
			planned_minutes,
			status,
			actual_start_at,
			current_agenda_id,
			COALESCE(decisions, ''),
			COALESCE(todo, '')
		FROM meetings
		WHERE 
			id = $1
			AND status = 'in_progress'
	`, id).Scan(
		&meeting.ID,
		&meeting.Title,
		&meeting.Description,
		&meeting.PlannedMinutes,
		&meeting.Status,
		&meeting.ActualStartAt,
		&meeting.CurrentAgendaID,
		&meeting.Decisions,
		&meeting.Todo,
	)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT
			id,
			title,
			COALESCE(purpose, ''),
			COALESCE(discussion_points, ''),
			COALESCE(questions, ''),
			COALESCE(memo, ''),
			planned_minutes,
			sort_order,
			actual_start_at,
			actual_end_at,
			elapsed_seconds
		FROM agendas
		WHERE meeting_id = $1
		ORDER BY sort_order
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meeting.Agendas = make([]model.AgendaSessionDetail, 0)

	for rows.Next() {
		var agenda model.AgendaSessionDetail

		if err := rows.Scan(
			&agenda.ID,
			&agenda.Title,
			&agenda.Purpose,
			&agenda.DiscussionPoints,
			&agenda.Questions,
			&agenda.Memo,
			&agenda.PlannedMinutes,
			&agenda.SortOrder,
			&agenda.ActualStartAt,
			&agenda.ActualEndAt,
			&agenda.ElapsedSeconds,
		); err != nil {
			return nil, err
		}

		meeting.Agendas = append(meeting.Agendas, agenda)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &meeting, nil
}

// ============================================
// 会議中：議題を戻す/進める
// ============================================
func ChangeCurrentAgenda(
	db *sql.DB,
	meetingID int64,
	req model.ChangeCurrentAgendaRequest,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 会議の現在議題を取得してロック
	var currentAgendaID *int64

	err = tx.QueryRow(`
		SELECT current_agenda_id
		FROM meetings
		WHERE id = $1
		FOR UPDATE
	`,
		meetingID,
	).Scan(&currentAgendaID)

	if err != nil {
		return err
	}

	// 移動先Agendaが対象会議に属しているか確認
	var targetAgendaID int64

	err = tx.QueryRow(`
		SELECT id
		FROM agendas
		WHERE
			id = $1
			AND meeting_id = $2
	`,
		req.TargetAgendaID,
		meetingID,
	).Scan(&targetAgendaID)

	if err != nil {
		return err
	}

	// 同じ議題への切り替えなら何もしない
	if currentAgendaID != nil &&
		*currentAgendaID == targetAgendaID {
		return nil
	}

	now := time.Now()

	// 現在議題の経過時間を累積して終了
	if currentAgendaID != nil {
		_, err = tx.Exec(`
			UPDATE agendas
			SET
				elapsed_seconds =
					elapsed_seconds +
					CASE
						WHEN actual_start_at IS NOT NULL
						THEN EXTRACT(
							EPOCH FROM ($1 - actual_start_at)
						)::INTEGER
						ELSE 0
					END,
				actual_end_at = $1,
				updated_at = NOW()
			WHERE
				id = $2
				AND meeting_id = $3
		`,
			now,
			*currentAgendaID,
			meetingID,
		)
		if err != nil {
			return err
		}
	}

	// 移動先議題を開始
	_, err = tx.Exec(`
		UPDATE agendas
		SET
			actual_start_at = $1,
			actual_end_at = NULL,
			updated_at = NOW()
		WHERE
			id = $2
			AND meeting_id = $3
	`,
		now,
		targetAgendaID,
		meetingID,
	)
	if err != nil {
		return err
	}

	// 会議の現在議題を変更
	_, err = tx.Exec(`
		UPDATE meetings
		SET
			current_agenda_id = $1,
			updated_at = NOW()
		WHERE id = $2
	`,
		targetAgendaID,
		meetingID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ============================================
// 会議中：一時保存
// ============================================
func SaveMeetingSession(
	db *sql.DB,
	id int64,
	req model.SaveMeetingSessionRequest,
) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 決定事項・TODOを保存
	result, err := tx.Exec(`
		UPDATE meetings
		SET
			decisions = $1,
			todo = $2,
			updated_at = NOW()
		WHERE
			id = $3
			AND status = 'in_progress'
	`,
		req.Decisions,
		req.Todo,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf(
			"meeting not found or session is not in progress",
		)
	}

	// 各Agendaのメモを保存
	for _, agenda := range req.Agendas {
		result, err := tx.Exec(`
			UPDATE agendas
			SET
				memo = $1,
				updated_at = NOW()
			WHERE
				id = $2
				AND meeting_id = $3
		`,
			agenda.Memo,
			agenda.ID,
			id,
		)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return fmt.Errorf(
				"agenda not found: id=%d",
				agenda.ID,
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}