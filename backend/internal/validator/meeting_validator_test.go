package validator

import (
	"strings"
	"testing"

	"meeting-pilot/internal/model"
)

func TestValidateCreateMeeting(t *testing.T) {
	tests := []struct {
		name        string
		req         model.CreateMeetingRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "正常なリクエスト",
			req: model.CreateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas: []model.CreateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "会議名が空",
			req: model.CreateMeetingRequest{
				Title:          "   ",
				PlannedMinutes: 60,
				Agendas: []model.CreateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "会議名は必須です",
		},
		{
			name: "会議時間が0分",
			req: model.CreateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 0,
				Agendas: []model.CreateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "会議時間は1分以上",
		},
		{
			name: "アジェンダが0件",
			req: model.CreateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas:        []model.CreateAgendaRequest{},
			},
			wantErr:     true,
			errContains: "アジェンダを1件以上",
		},
		{
			name: "2件目のアジェンダ名が空",
			req: model.CreateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas: []model.CreateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
					{
						Title:          "",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "議題2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateMeeting(tt.req)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf(
						"エラーを期待していませんが発生しました: %v",
						err,
					)
				}

				return
			}

			if err == nil {
				t.Fatal("エラーを期待しましたがnilでした")
			}

			if !strings.Contains(
				err.Error(),
				tt.errContains,
			) {
				t.Errorf(
					"エラーに%qが含まれていません: %v",
					tt.errContains,
					err,
				)
			}
		})
	}
}

func TestValidateUpdateMeeting(t *testing.T) {
	tests := []struct {
		name        string
		req         model.UpdateMeetingRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "正常なリクエスト",
			req: model.UpdateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas: []model.UpdateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "会議名が空",
			req: model.UpdateMeetingRequest{
				Title:          "   ",
				PlannedMinutes: 60,
				Agendas: []model.UpdateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "会議名は必須です",
		},
		{
			name: "会議時間が0分",
			req: model.UpdateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 0,
				Agendas: []model.UpdateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "会議時間は1分以上",
		},
		{
			name: "アジェンダが0件",
			req: model.UpdateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas:        []model.UpdateAgendaRequest{},
			},
			wantErr:     true,
			errContains: "アジェンダを1件以上",
		},
		{
			name: "2件目のアジェンダ名が空",
			req: model.UpdateMeetingRequest{
				Title:          "週次会議",
				PlannedMinutes: 60,
				Agendas: []model.UpdateAgendaRequest{
					{
						Title:          "進捗確認",
						PlannedMinutes: 30,
					},
					{
						Title:          "",
						PlannedMinutes: 30,
					},
				},
			},
			wantErr:     true,
			errContains: "議題2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateMeeting(tt.req)

			if !tt.wantErr {
				if err != nil {
					t.Fatalf(
						"エラーを期待していませんが発生しました: %v",
						err,
					)
				}

				return
			}

			if err == nil {
				t.Fatal("エラーを期待しましたがnilでした")
			}

			if !strings.Contains(
				err.Error(),
				tt.errContains,
			) {
				t.Errorf(
					"エラーに%qが含まれていません: %v",
					tt.errContains,
					err,
				)
			}
		})
	}
}