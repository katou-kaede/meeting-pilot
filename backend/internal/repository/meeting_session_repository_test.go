package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"meeting-pilot/internal/model"
)

func TestChangeCurrentAgendaSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	const (
		meetingID      int64 = 1
		currentAgendaID int64 = 10
		targetAgendaID  int64 = 20
	)

	req := model.ChangeCurrentAgendaRequest{
		TargetAgendaID: targetAgendaID,
	}

	// トランザクション開始
	mock.ExpectBegin()

	// 現在議題を取得
	mock.ExpectQuery(
		`SELECT current_agenda_id\s+FROM meetings\s+WHERE id = \$1\s+FOR UPDATE`,
	).
		WithArgs(meetingID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"current_agenda_id"},
			).AddRow(currentAgendaID),
		)

	// 移動先議題を確認
	mock.ExpectQuery(
		`SELECT id\s+FROM agendas\s+WHERE\s+id = \$1\s+AND meeting_id = \$2`,
	).
		WithArgs(
			targetAgendaID,
			meetingID,
		).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id"},
			).AddRow(targetAgendaID),
		)

	// 現在議題を終了
	mock.ExpectExec(
		`UPDATE agendas\s+SET\s+elapsed_seconds`,
	).
		WithArgs(
			sqlmock.AnyArg(),
			currentAgendaID,
			meetingID,
		).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	// 移動先議題を開始
	mock.ExpectExec(
		`UPDATE agendas\s+SET\s+actual_start_at`,
	).
		WithArgs(
			sqlmock.AnyArg(),
			targetAgendaID,
			meetingID,
		).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	// 会議の現在議題を更新
	mock.ExpectExec(
		`UPDATE meetings\s+SET\s+current_agenda_id`,
	).
		WithArgs(
			targetAgendaID,
			meetingID,
		).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	// Commit
	mock.ExpectCommit()

	err = ChangeCurrentAgenda(
		db,
		meetingID,
		req,
	)
	if err != nil {
		t.Fatalf(
			"エラーを期待していませんが発生しました: %v",
			err,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestChangeCurrentAgendaTargetNotFound(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	const (
		meetingID      int64 = 1
		currentAgendaID int64 = 10
		targetAgendaID  int64 = 99
	)

	req := model.ChangeCurrentAgendaRequest{
		TargetAgendaID: targetAgendaID,
	}

	mock.ExpectBegin()

	mock.ExpectQuery(
		`SELECT current_agenda_id\s+FROM meetings\s+WHERE id = \$1\s+FOR UPDATE`,
	).
		WithArgs(meetingID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"current_agenda_id"},
			).AddRow(currentAgendaID),
		)

	mock.ExpectQuery(
		`SELECT id\s+FROM agendas\s+WHERE\s+id = \$1\s+AND meeting_id = \$2`,
	).
		WithArgs(
			targetAgendaID,
			meetingID,
		).
		WillReturnError(
			errors.New("target agenda not found"),
		)

	// defer tx.Rollback() が実行される
	mock.ExpectRollback()

	err = ChangeCurrentAgenda(
		db,
		meetingID,
		req,
	)

	if err == nil {
		t.Fatal(
			"エラーを期待しましたがnilでした",
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestChangeCurrentAgendaCurrentAgendaUpdateError(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"sqlmockの作成に失敗しました: %v",
			err,
		)
	}
	defer db.Close()

	const (
		meetingID      int64 = 1
		currentAgendaID int64 = 10
		targetAgendaID  int64 = 20
	)

	req := model.ChangeCurrentAgendaRequest{
		TargetAgendaID: targetAgendaID,
	}

	databaseError := errors.New(
		"current agenda update failed",
	)

	mock.ExpectBegin()

	mock.ExpectQuery(
		`SELECT current_agenda_id\s+FROM meetings\s+WHERE id = \$1\s+FOR UPDATE`,
	).
		WithArgs(meetingID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"current_agenda_id"},
			).AddRow(currentAgendaID),
		)

	mock.ExpectQuery(
		`SELECT id\s+FROM agendas\s+WHERE\s+id = \$1\s+AND meeting_id = \$2`,
	).
		WithArgs(
			targetAgendaID,
			meetingID,
		).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id"},
			).AddRow(targetAgendaID),
		)

	mock.ExpectExec(
		`UPDATE agendas\s+SET\s+elapsed_seconds`,
	).
		WithArgs(
			sqlmock.AnyArg(),
			currentAgendaID,
			meetingID,
		).
		WillReturnError(databaseError)

	mock.ExpectRollback()

	err = ChangeCurrentAgenda(
		db,
		meetingID,
		req,
	)

	if !errors.Is(err, databaseError) {
		t.Errorf(
			"error = %v, want %v",
			err,
			databaseError,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(
			"未実行のSQL期待値があります: %v",
			err,
		)
	}
}

func TestPauseMeeting(t *testing.T) {
	databaseError := errors.New(
		"agenda update failed",
	)

	tests := []struct {
		name               string
		meetingRowsAffected int64
		agendaUpdateError   error
		wantErr             bool
		wantCommit          bool
	}{
		{
			name:                "正常に一時停止できる",
			meetingRowsAffected: 1,
			wantErr:             false,
			wantCommit:          true,
		},
		{
			name:                "停止可能な会議が存在しない",
			meetingRowsAffected: 0,
			wantErr:             true,
			wantCommit:          false,
		},
		{
			name:                "議題更新に失敗したらRollback",
			meetingRowsAffected: 1,
			agendaUpdateError:   databaseError,
			wantErr:             true,
			wantCommit:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			const meetingID int64 = 1

			mock.ExpectBegin()

			mock.ExpectExec(
				`UPDATE meetings\s+SET\s+paused_at = NOW\(\)`,
			).
				WithArgs(meetingID).
				WillReturnResult(
					sqlmock.NewResult(
						0,
						tt.meetingRowsAffected,
					),
				)

			// 会議更新が成功した場合のみ議題更新へ進む
			if tt.meetingRowsAffected > 0 {
				agendaExpectation := mock.ExpectExec(
					`UPDATE agendas\s+SET\s+elapsed_seconds`,
				).
					WithArgs(meetingID)

				if tt.agendaUpdateError != nil {
					agendaExpectation.WillReturnError(
						tt.agendaUpdateError,
					)
				} else {
					agendaExpectation.WillReturnResult(
						sqlmock.NewResult(0, 1),
					)
				}
			}

			if tt.wantCommit {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err = PauseMeeting(db, meetingID)

			if tt.wantErr {
				if err == nil {
					t.Fatal(
						"エラーを期待しましたがnilでした",
					)
				}
			} else if err != nil {
				t.Fatalf(
					"エラーを期待していませんが発生しました: %v",
					err,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}

func TestResumeMeeting(t *testing.T) {
	databaseError := errors.New(
		"agenda resume failed",
	)

	tests := []struct {
		name                string
		meetingRowsAffected int64
		agendaUpdateError   error
		wantErr             bool
		wantCommit          bool
	}{
		{
			name:                "正常に会議を再開できる",
			meetingRowsAffected: 1,
			wantErr:             false,
			wantCommit:          true,
		},
		{
			name:                "再開可能な会議が存在しない",
			meetingRowsAffected: 0,
			wantErr:             true,
			wantCommit:          false,
		},
		{
			name:                "議題の再開に失敗したらRollback",
			meetingRowsAffected: 1,
			agendaUpdateError:   databaseError,
			wantErr:             true,
			wantCommit:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			const meetingID int64 = 1

			mock.ExpectBegin()

			mock.ExpectExec(
				`UPDATE meetings\s+SET\s+total_paused_seconds`,
			).
				WithArgs(meetingID).
				WillReturnResult(
					sqlmock.NewResult(
						0,
						tt.meetingRowsAffected,
					),
				)

			if tt.meetingRowsAffected > 0 {
				agendaExpectation := mock.ExpectExec(
					`UPDATE agendas\s+SET\s+actual_start_at = NOW\(\)`,
				).
					WithArgs(meetingID)

				if tt.agendaUpdateError != nil {
					agendaExpectation.WillReturnError(
						tt.agendaUpdateError,
					)
				} else {
					agendaExpectation.WillReturnResult(
						sqlmock.NewResult(0, 1),
					)
				}
			}

			if tt.wantCommit {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err = ResumeMeeting(db, meetingID)

			if tt.wantErr {
				if err == nil {
					t.Fatal(
						"エラーを期待しましたがnilでした",
					)
				}
			} else if err != nil {
				t.Fatalf(
					"エラーを期待していませんが発生しました: %v",
					err,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}

func TestCompleteMeeting(t *testing.T) {
	databaseError := errors.New("agenda update failed")

	tests := []struct {
		name                string
		meetingRowsAffected int64
		agendaUpdateError   error
		wantErr             bool
		wantCommit          bool
	}{
		{
			name:                "正常に会議を完了できる",
			meetingRowsAffected: 1,
			wantErr:             false,
			wantCommit:          true,
		},
		{
			name:                "完了可能な会議が存在しない",
			meetingRowsAffected: 0,
			wantErr:             true,
			wantCommit:          false,
		},
		{
			name:                "Agenda更新に失敗したらRollback",
			meetingRowsAffected: 1,
			agendaUpdateError:   databaseError,
			wantErr:             true,
			wantCommit:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(
					"sqlmockの作成に失敗しました: %v",
					err,
				)
			}
			defer db.Close()

			const meetingID int64 = 1

			req := model.CompleteMeetingRequest{
				Decisions: "新機能を実装する",
				Todo:      "次回までに仕様を確認する",
				Agendas: []model.CompleteAgendaRequest{
					{
						ID:   10,
						Memo: "議題1のメモ",
					},
					{
						ID:   20,
						Memo: "議題2のメモ",
					},
				},
			}

			mock.ExpectBegin()

			// 会議結果・ステータス更新
			mock.ExpectExec(
				`UPDATE meetings\s+SET\s+status = \$1`,
			).
				WithArgs(
					"completed",
					req.Decisions,
					req.Todo,
					meetingID,
				).
				WillReturnResult(
					sqlmock.NewResult(
						0,
						tt.meetingRowsAffected,
					),
				)

			// 会議更新に成功した場合のみAgenda更新へ進む
			if tt.meetingRowsAffected > 0 {
				for index, agenda := range req.Agendas {
					expectation := mock.ExpectExec(
						`UPDATE agendas\s+SET\s+memo = \$1`,
					).
						WithArgs(
							agenda.Memo,
							agenda.ID,
							meetingID,
						)

					// Agenda更新エラーは1件目で発生させる
					if index == 0 &&
						tt.agendaUpdateError != nil {
						expectation.WillReturnError(
							tt.agendaUpdateError,
						)

						break
					}

					expectation.WillReturnResult(
						sqlmock.NewResult(0, 1),
					)
				}
			}

			if tt.wantCommit {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err = CompleteMeeting(
				db,
				meetingID,
				req,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal(
						"エラーを期待しましたがnilでした",
					)
				}
			} else if err != nil {
				t.Fatalf(
					"エラーを期待していませんが発生しました: %v",
					err,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf(
					"未実行のSQL期待値があります: %v",
					err,
				)
			}
		})
	}
}