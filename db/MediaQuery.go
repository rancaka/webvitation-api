package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rancaka/webvitation-api/models"
)

type MediaType string

const MediaTypeImage MediaType = "image"
const MediaTypeVideo MediaType = "video"

const mediaColumns = `
	media_id,
	coalesce(path, '') as path,
	coalesce(video_poster_path, '') as video_poster_path,
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

func insertMedia(ctx context.Context, eventID, path, videoPosterPath string, mediaType MediaType) error {

	_, err := db.ExecContext(
		ctx,
		`INSERT
			INTO event_media (event_id, path, video_poster_path, type)
			VALUES (uuid_to_bin(?), ?, ?, ?)`,
		eventID,
		path,
		videoPosterPath,
		mediaType,
	)

	return err
}
