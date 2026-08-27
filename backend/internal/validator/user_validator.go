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

	for _, r := range req.Password {
		if r < 0x21 || r > 0x7E {
			return errors.New(
				"パスワードは半角英数字・記号で入力してください",
			)
		}
	}

	if len(req.Password) < 8 {
		return errors.New(
			"パスワードは8文字以上で入力してください",
		)
	}

	if len(req.Password) > 72 {
		return errors.New(
			"パスワードは72文字以内で入力してください",
		)
	}

	return nil
}
