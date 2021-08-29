package models

type Media struct {
	MediaID         string `db:"media_id"`
	Path            string `db:"path"`
	VideoPosterPath string `db:"video_poster_path"`
	EventID         string `db:"event_id"`
	MediaType       string `db:"type"`
}
