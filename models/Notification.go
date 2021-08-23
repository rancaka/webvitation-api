package models

import (
	"errors"
	"time"
)

type NotificationType string

const (
	NotificationCheckIn NotificationType = "CHECK_IN"
)

type Notification struct {
	NotificationID string                 `db:"notification_id"`
	EventID        string                 `db:"event_id"`
	Type           NotificationType       `db:"type"`
	Data           map[string]interface{} `db:"data"`
	CreatedAt      time.Time              `db:"created_at"`

	CheckIn *CheckIn
}

func (n *Notification) SetData() error {
	var ok = true

	switch n.Type {
	case NotificationCheckIn:

		checkIn := new(CheckIn)
		checkIn.CreatorUID, ok = n.Data["creator_uid"].(string)
		if !ok {
			return errors.New("missing creator_uid")
		}

		checkIn.InviteeID, ok = n.Data["invitee_id"].(string)
		if !ok {
			return errors.New("missing invitee_id")
		}

		checkIn.InviteeName, ok = n.Data["invitee_name"].(string)
		if !ok {
			return errors.New("missing invitee_name")
		}

		checkIn.CheckedInAt, ok = n.Data["checked_in_at"].(string)
		if !ok {
			return errors.New("missing checked_in_at")
		}

		checkIn.CheckedInAtUnix, ok = n.Data["checked_in_at_unix"].(int64)
		if !ok {
			return errors.New("missing checked_in_at_unix")
		}
	}

	n.Data = nil
	return nil
}

type CheckIn struct {
	CreatorUID      string
	InviteeID       string
	InviteeName     string
	CheckedInAt     string
	CheckedInAtUnix int64
}
