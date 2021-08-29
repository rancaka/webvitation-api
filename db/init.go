package db

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rancaka/webvitation-api/models"
	"github.com/rancaka/webvitation-api/util"
)

var (
	db *sqlx.DB = util.GetDB()

	// EventQuery
	GetEventsByUID    func(ctx context.Context, UID string) ([]*models.Event, error)               = getEventsByUID
	GetEventByEventID func(ctx context.Context, eventID string) (*models.Event, error)             = getEventByEventID
	GetEventBySlug    func(ctx context.Context, slug string) (*models.Event, error)                = getEventBySlug
	CreateEvent       func(ctx context.Context, eventRequest *models.EventRequest) (string, error) = createEvent

	// InviteeQuery
	GetInviteesByEventID            func(ctx context.Context, eventID string) ([]*models.Invitee, error)                                                                        = getInviteesByEventID
	GetInviteeByEventIDAndSlug      func(ctx context.Context, eventID, inviteeSlug string) (*models.Invitee, error)                                                             = getInviteeByEventIDAndSlug
	GetInviteeByEventIDAndInviteeID func(ctx context.Context, eventID, inviteeID string) (*models.Invitee, error)                                                               = getInviteeByEventIDAndInviteeID
	CreateInviteesBulk              func(ctx context.Context, eventID string, inviteeRequests []*models.InviteeRequest) (duplicateInvitees []*models.InviteeRequest, err error) = createInviteesBulk
	CreateInviteesBulk2             func(ctx context.Context, eventID string, inviteeRequests []*models.InviteeRequest) error                                                   = createInviteesBulk2

	// NotificationQuery
	GetNotifications   func(ctx context.Context, eventID string) ([]*models.Notification, error)    = getNotificationsByEventID
	CreateNotification func(ctx context.Context, notification *models.Notification) (string, error) = createNotification

	// MediaQuery
	GetMediaByEventID func(ctx context.Context, eventID string) ([]*models.Media, error)                          = getMediaByEventID
	InsertMedia       func(ctx context.Context, eventID, path, videoPosterPath string, mediaType MediaType) error = insertMedia
)

func getLastInsertedUUID(row *sqlx.Row) (string, error) {

	var lastInsertedUUID string = ""
	err := row.Scan(&lastInsertedUUID)
	if err != nil {
		log.Printf("[ERROR] db.getLastInsertedUUID -> row.Scan: %v\n", err)
		return "", err
	}

	return lastInsertedUUID, nil
}
