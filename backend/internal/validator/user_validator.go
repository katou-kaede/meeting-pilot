package validator

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"

	"meeting-pilot/internal/model"
)

func ValidateCreateUser(
	req model.CreateUserRequest,
) error {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" {
		return errors.New("氏名は必須です")
	}

	if utf8.RuneCountInString(name) > 100 {
		return errors.New("氏名は100文字以内で入力してください")
	}

	if email == "" {
		return errors.New("メールアドレスは必須です")
	}

	if len(email) > 255 {
		return errors.New("メールアドレスは255文字以内で入力してください")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("メールアドレスの形式が正しくありません")
	}

	if len(req.Password) < 8 {
		return errors.New("パスワードは8文字以上で入力してください")
	}

	// bcryptは72バイトを超えるパスワードを受け付けない
	if len([]byte(req.Password)) > 72 {
		return errors.New("パスワードが長すぎます")
	}

	return nil
}
