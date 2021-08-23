package models

import "database/sql"

type Event struct {
	EventID           string         `db:"event_id"`
	Slug              string         `db:"slug"`
	InvitationMessage sql.NullString `db:"invitation_message"`
	CreatorUID        string         `db:"creator_uid"`
	Timestamp

	Invitees      []*Invitee
	Notifications []*Notification
}

type EventRequest struct {
	Slug              string `binding:"required"`
	InvitationMessage *string
	CreatorUID        string

	Invitees []*InviteeRequest
}
