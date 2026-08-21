package validator

import (
	"errors"
	"strings"

	"meeting-pilot/internal/model"
)

func ValidateCreateMeeting(
	req model.CreateMeetingRequest,
) error {

	if strings.TrimSpace(req.Title) == "" {
		return errors.New("会議名は必須です")
	}

	if req.PlannedMinutes <= 0 {
		return errors.New("会議時間は1分以上入力してください")
	}

	if len(req.Agendas) == 0 {
		return errors.New("アジェンダを1件以上入力してください")
	}

	for _, agenda := range req.Agendas {
		if strings.TrimSpace(agenda.Title) == "" {
			return errors.New("アジェンダ名は必須です")
		}

		if agenda.PlannedMinutes <= 0 {
			return errors.New("アジェンダ時間は1分以上入力してください")
		}
	}

	return nil
}

func ValidateUpdateMeeting(
	req model.UpdateMeetingRequest,
) error {

	if strings.TrimSpace(req.Title) == "" {
		return errors.New("会議名は必須です")
	}

	if req.PlannedMinutes <= 0 {
		return errors.New("会議時間は1分以上入力してください")
	}

	if len(req.Agendas) == 0 {
		return errors.New("アジェンダを1件以上入力してください")
	}

	for _, agenda := range req.Agendas {
		if strings.TrimSpace(agenda.Title) == "" {
			return errors.New("アジェンダ名は必須です")
		}

		if agenda.PlannedMinutes <= 0 {
			return errors.New("アジェンダ時間は1分以上入力してください")
		}
	}

	return nil
}
