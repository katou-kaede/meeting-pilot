package repository

import (
	"database/sql"
	"fmt"
	"log"
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
		log.Printf("insert meetings failed: %v", err)
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
			log.Printf("insert agendas failed: %v", err)
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
		log.Println(err)
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
				log.Printf("agenda delete error: %v", err)
				return err
			}
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
			log.Printf("agenda update error: %v", err)
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
			log.Printf("agenda insert error: %v", err)
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
	_, err = tx.Exec(`
		UPDATE meetings
		SET
			status = $1,
			decisions = $2,
			todo = $3,
			actual_end_at = NOW(),
			updated_at = NOW()
		WHERE id = $4
	`,
		"completed",
		req.Decisions,
		req.Todo,
		id,
	)
	if err != nil {
		return err
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
		WHERE id = $1
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