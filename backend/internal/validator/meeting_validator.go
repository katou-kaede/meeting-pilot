package validator

import (
	"errors"
	"strings"
	"fmt"
	"unicode/utf8"

	"meeting-pilot/internal/model"
)

func ValidateCreateMeeting(
	req model.CreateMeetingRequest,
) error {
	title := strings.TrimSpace(req.Title)

	if title == "" {
		return errors.New("会議名は必須です")
	}

	if utf8.RuneCountInString(title) > 200 {
		return errors.New("会議名は200文字以内で入力してください")
	}

	if utf8.RuneCountInString(req.TargetName) > 200 {
		return errors.New("会議相手は200文字以内で入力してください")
	}

	if req.PlannedMinutes <= 0 {
		return errors.New("会議時間は1分以上入力してください")
	}

	if len(req.Agendas) == 0 {
		return errors.New("アジェンダを1件以上入力してください")
	}

	for index, agenda := range req.Agendas {
		agendaTitle := strings.TrimSpace(agenda.Title)

		if agendaTitle == "" {
			return fmt.Errorf(
				"議題%dのアジェンダ名は必須です",
				index+1,
			)
		}

		if utf8.RuneCountInString(agendaTitle) > 200 {
			return fmt.Errorf(
				"議題%dのアジェンダ名は200文字以内で入力してください",
				index+1,
			)
		}

		if agenda.PlannedMinutes <= 0 {
			return fmt.Errorf(
				"議題%dの時間は1分以上入力してください",
				index+1,
			)
		}
	}

	return nil
}

func ValidateUpdateMeeting(
	req model.UpdateMeetingRequest,
) error {

	title := strings.TrimSpace(req.Title)

	if title == "" {
		return errors.New("会議名は必須です")
	}

	if utf8.RuneCountInString(title) > 200 {
		return errors.New("会議名は200文字以内で入力してください")
	}

	if utf8.RuneCountInString(req.TargetName) > 200 {
		return errors.New("会議相手は200文字以内で入力してください")
	}

	if req.PlannedMinutes <= 0 {
		return errors.New("会議時間は1分以上入力してください")
	}

	if len(req.Agendas) == 0 {
		return errors.New("アジェンダを1件以上入力してください")
	}

	for index, agenda := range req.Agendas {
		agendaTitle := strings.TrimSpace(agenda.Title)

		if agendaTitle == "" {
			return fmt.Errorf(
				"議題%dのアジェンダ名は必須です",
				index+1,
			)
		}

		if utf8.RuneCountInString(agendaTitle) > 200 {
			return fmt.Errorf(
				"議題%dのアジェンダ名は200文字以内で入力してください",
				index+1,
			)
		}

		if agenda.PlannedMinutes <= 0 {
			return fmt.Errorf(
				"議題%dの時間は1分以上入力してください",
				index+1,
			)
		}
	}

	return nil
}
