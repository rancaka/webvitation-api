package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rancaka/webvitation-api/models"
)

const eventColumns = `
	event_id,
	slug,
	invitation_message,
	creator_uid,
	created_at,
	updated_at
`

func getEventsByUID(ctx context.Context, creatorUID string) ([]*models.Event, error) {

	events := []*models.Event{}
	query := "SELECT " + eventColumns + " FROM events WHERE creator_uid = ? ORDER BY updated_at DESC"
	err := sqlx.SelectContext(ctx, db, &events, query, creatorUID)
	if err != nil {
		log.Printf("[ERROR] db.GetEventsByUID -> sqlx.SelectContext: (%v)\n", err)
		return nil, err
	}

	return events, nil
}

func getEventByEventID(ctx context.Context, eventID string) (*models.Event, error) {

	event := new(models.Event)
	query := "SELECT " + eventColumns + " FROM events WHERE event_id = uuid_to_bin(?) LIMIT 1"
	err := sqlx.GetContext(ctx, db, event, query, eventID)
	if err != nil {
		log.Printf("[ERROR] db.GetEventByEventID -> sqlx.GetContext: (%v)\n", err)
		return nil, err
	}

	return event, nil
}

func getEventBySlug(ctx context.Context, slug string) (*models.Event, error) {

	event := new(models.Event)
	query := "SELECT " + eventColumns + " FROM events WHERE slug = ? LIMIT 1"
	err := sqlx.GetContext(ctx, db, event, query, slug)
	if err != nil {
		log.Printf("[ERROR] db.GetEventBySlug -> sqlx.GetContext: (%v)\n", err)
		return nil, err
	}

	return event, nil
}

func createEvent(ctx context.Context, eventRequest *models.EventRequest) (string, error) {
	row := db.QueryRowxContext(
		ctx,
		`INSERT INTO events (slug, invitation_message, creator_uid) VALUES (?, ?, ?) RETURNING event_id`,
		eventRequest.Slug,
		eventRequest.InvitationMessage,
		eventRequest.CreatorUID,
	)

	return getLastInsertedUUID(row)
}
