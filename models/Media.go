package models

type Media struct {
	MediaID        string `db:"media_id"`
	URL            string `db:"url"`
	EventID        string `db:"event_id"`
	MediaType      string `db:"type"`
	VideoPosterURL string `db:"video_poster_url"`
}
