package models

type Invitee struct {
	InviteeID string `db:"invitee_id"`
	Name      string `db:"name"`
	Slug      string `db:"slug"`
	EventID   string `db:"event_id"`
	Timestamp
}

type InviteeRequest struct {
	Name    string
	Slug    string
	EventID string `db:"event_id"`
}
