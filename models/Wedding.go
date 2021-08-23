package models

import "time"

type Wedding struct {
	WeddingID string `db:"wedding_id"`

	MetaDescription      string `db:"meta_description"`
	OGPhotoURL           string `db:"og_photo_url"`
	CoverPhotoURL        string `db:"cover_photo_url"`
	MainQuote            string `db:"main_quote"`
	BacksoundURL         string `db:"backsound_url"`
	CheckInBackgroundURL string `db:"check_in_background_url"`

	BrideNickName     string `db:"bride_nick_name"`
	BrideFullName     string `db:"bride_full_name"`
	BrideParentsLabel string `db:"bride_parents_label"`
	BridePhotoURL     string `db:"bride_photo_url"`
	BrideInstagram    string `db:"bride_instagram"`

	GroomNickName     string `db:"groom_nick_name"`
	GroomFullName     string `db:"groom_full_name"`
	GroomParentsLabel string `db:"groom_parents_label"`
	GroomPhotoURL     string `db:"groom_photo_url"`
	GroomInstagram    string `db:"groom_instagram"`

	CeremonyDate         time.Time `db:"ceremony_date"`
	CeremonyTimeFrom     time.Time `db:"ceremony_time_from"`
	CeremonyTimeTo       time.Time `db:"ceremony_time_to"`
	CeremonyLocationURL  string    `db:"ceremony_location_url"`
	CeremonyLocationName string    `db:"ceremony_location_name"`

	ReceptionDate         time.Time `db:"reception_date"`
	ReceptionTimeFrom     time.Time `db:"reception_time_from"`
	ReceptionTimeTo       time.Time `db:"reception_time_to"`
	ReceptionLocationURL  string    `db:"reception_location_url"`
	ReceptionLocationName string    `db:"reception_location_name"`

	EventID    string `db:"event_id"`
	TemplateID string `db:"template_id"`

	Timestamp
}
