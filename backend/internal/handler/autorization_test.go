package handler

import "testing"

func TestHasMeetingSessionEditPermission(t *testing.T) {
	editorUserID := int64(2)

	tests := []struct {
		name         string
		role         string
		userID       int64
		editorUserID *int64
		want         bool
	}{
		{
			name:         "editor未設定ならownerは編集可能",
			role:         "owner",
			userID:       1,
			editorUserID: nil,
			want:         true,
		},
		{
			name:         "editor未設定でもviewerは編集不可",
			role:         "viewer",
			userID:       2,
			editorUserID: nil,
			want:         false,
		},
		{
			name:         "editor未設定でもeditorロールは編集不可",
			role:         "editor",
			userID:       2,
			editorUserID: nil,
			want:         false,
		},
		{
			name:         "editor設定済みならeditor本人は編集可能",
			role:         "editor",
			userID:       2,
			editorUserID: &editorUserID,
			want:         true,
		},
		{
			name:         "editor設定済みならownerは編集不可",
			role:         "owner",
			userID:       1,
			editorUserID: &editorUserID,
			want:         false,
		},
		{
			name:         "editor設定済みなら別のviewerは編集不可",
			role:         "viewer",
			userID:       3,
			editorUserID: &editorUserID,
			want:         false,
		},
		{
			name:         "editorロールでも設定されたIDと異なれば編集不可",
			role:         "editor",
			userID:       3,
			editorUserID: &editorUserID,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasMeetingSessionEditPermission(
				tt.role,
				tt.userID,
				tt.editorUserID,
			)

			if got != tt.want {
				t.Errorf(
					"hasMeetingSessionEditPermission() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}