package db

import (
	"context"
	"log"

	"github.com/rancaka/webvitation-web/models"
)

const (
	CheckInAtLayout     = "02/01/2006 15:04:05 PM"
	notificationColumns = `
		notification_id,
		type,
		data,
		event_id,
		created_at
	`
)

func getNotificationsByEventID(ctx context.Context, eventID string) ([]*models.Notification, error) {

	notifications := []*models.Notification{}
	query := "SELECT " + notificationColumns + " FROM event_notifications WHERE event_id = ?"

	rows, err := db.QueryxContext(ctx, query, eventID)
	if err != nil {
		log.Printf("[ERROR] db.getNotificationsByEventID -> db.QueryxContext: (%v)\n", err)
		return nil, err
	}

	notification := new(models.Notification)
	for rows.Next() {
		err = rows.StructScan(notification)
		if err != nil {
			log.Printf("[ERROR] db.getNotificationsByEventID -> rows.StructScan: (%v)\n", err)
			return nil, err
		}

		err = notification.SetData()
		if err != nil {
			log.Printf("[ERROR] db.getNotificationsByEventID -> notification.SetData(): (%v)\n", err)
			return nil, err
		}
	}

	return notifications, nil
}

func createNotification(ctx context.Context, notification *models.Notification) (string, error) {
	// TODO: implementation
	return "", nil
}
