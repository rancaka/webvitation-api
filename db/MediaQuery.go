package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rancaka/webvitation-web/models"
)

const mediaColumns = `
	media_id,
	coalesce(url, '') as url,
	coalesce(video_poster_url, '') as video_poster_url,
	event_id,
	type
`

func getMediaByEventID(ctx context.Context, eventID string) ([]*models.Media, error) {

	media := []*models.Media{}
	query := "SELECT " + mediaColumns + " FROM event_media WHERE event_id = ?"
	err := sqlx.SelectContext(ctx, db, &media, query, eventID)
	if err != nil {
		log.Printf("[ERROR] db.GetMediaByEventID -> sqlx.SelectContext: (%v)\n", err)
		return nil, err
	}

	return media, nil
}
