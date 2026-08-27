package validator

import (
	"strings"
	"testing"

	"meeting-pilot/internal/model"
)

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		req         model.CreateUserRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "正常なリクエスト",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "氏名の前後に空白があっても正常",
			req: model.CreateUserRequest{
				Name:     "  Employee1  ",
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "メールアドレスに大文字が含まれていても正常",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "Employee1@Example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "氏名が空",
			req: model.CreateUserRequest{
				Name:     "",
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "氏名は必須です",
		},
		{
			name: "氏名が空白のみ",
			req: model.CreateUserRequest{
				Name:     "   ",
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "氏名は必須です",
		},
		{
			name: "氏名が100文字を超えている",
			req: model.CreateUserRequest{
				Name:     strings.Repeat("あ", 101),
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "氏名は100文字以内で入力してください",
		},
		{
			name: "氏名が100文字ちょうど",
			req: model.CreateUserRequest{
				Name:     strings.Repeat("あ", 100),
				Email:    "employee1@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "メールアドレスが空",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "メールアドレスは必須です",
		},
		{
			name: "メールアドレスが空白のみ",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "   ",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "メールアドレスは必須です",
		},
		{
			name: "メールアドレスが255文字を超えている",
			req: model.CreateUserRequest{
				Name: "Employee1",
				Email: strings.Repeat("a", 244) +
					"@example.com",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "メールアドレスは255文字以内で入力してください",
		},
		{
			name: "メールアドレスの形式が不正",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr:     true,
			errContains: "メールアドレスの形式が正しくありません",
		},
		{
			name: "パスワードに日本語が含まれる",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "パスワード123",
			},
			wantErr:     true,
			errContains: "パスワードは半角英数字・記号で入力してください",
		},
		{
			name: "パスワードに半角空白が含まれる",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "pass 1234",
			},
			wantErr:     true,
			errContains: "パスワードは半角英数字・記号で入力してください",
		},
		{
			name: "パスワードが7文字",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "pass123",
			},
			wantErr:     true,
			errContains: "パスワードは8文字以上で入力してください",
		},
		{
			name: "パスワードが8文字ちょうど",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "pass1234",
			},
			wantErr: false,
		},
		{
			name: "パスワードに記号が含まれていても正常",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: "pass!234",
			},
			wantErr: false,
		},
		{
			name: "パスワードが72文字ちょうど",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: strings.Repeat("a", 72),
			},
			wantErr: false,
		},
		{
			name: "パスワードが73文字",
			req: model.CreateUserRequest{
				Name:     "Employee1",
				Email:    "employee1@example.com",
				Password: strings.Repeat("a", 73),
			},
			wantErr:     true,
			errContains: "パスワードは72文字以内で入力してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateUser(tt.req)

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
				t.Fatal(
					"エラーを期待しましたがnilでした",
				)
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