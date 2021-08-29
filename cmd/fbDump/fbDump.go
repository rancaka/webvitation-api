package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rancaka/webvitation-api/db"
	"github.com/rancaka/webvitation-api/models"
	"github.com/rancaka/webvitation-api/util"
)

var (
	fb *util.FirebaseApp
)

func init() {
	fb = util.GetFirebaseApp()
}

func main() {

	eventRequest := &models.EventRequest{
		Slug:       "melati-yogy-wedding",
		CreatorUID: "BglUnqnvRqWjsBlVOYJKC4RmB9f2",
	}

	newEventID, err := db.CreateEvent(context.Background(), eventRequest)
	if err != nil {
		log.Fatalln(err)
	}

	colRef := fb.Firestore.Collection("events/melati-yogy-wedding/guests")
	docSnapshots, err := colRef.Documents(context.Background()).GetAll()
	if err != nil {
		log.Fatalln(err)
		return
	}

	invitees := make([]models.InviteeRequest, len(docSnapshots))
	for index, snapshot := range docSnapshots {
		data := snapshot.Data()
		invitee := models.InviteeRequest{
			Name: data["name"].(string),
			Slug: snapshot.Ref.ID,
		}
		fmt.Printf("[%d] invitee: %+v\n", index, invitee)
		invitees[index] = invitee
	}

	err = db.CreateInviteesBulk2(context.Background(), newEventID, invitees)
	if err != nil {
		log.Fatalln(err)
	}

	_invitees, err := db.GetInviteesByEventID(context.Background(), newEventID)
	if err != nil {
		log.Fatalln(err)
	}

	colRef = fb.Firestore.Collection("events/melati-yogy-wedding/check-in")
	docSnapshots, err = colRef.Documents(context.Background()).GetAll()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for index, snapshot := range docSnapshots {

		data := snapshot.Data()
		name := data["name"].(string)
		checkedInAt := data["checkedInAt"].(string)
		createdAt, err := time.Parse("1/02/2006, 15:04:05 PM", checkedInAt)
		if err != nil {
			log.Fatalln(err)
		}

		var slug string
		value, ok := data["slug"].(string)
		if ok && value != "" {
			slug = value
		}

		fmt.Printf("[%d] check-in: %+v\n", index, data)

		checkInLogRequest := &models.CheckInLogRequest{
			EventID:     newEventID,
			InviteeID:   "",
			InviteeName: name,
			Timestamp: models.Timestamp{
				CreatedAt: createdAt,
			},
		}

		if slug == "manual" {
			_, err = db.CreateCheckInLog(context.Background(), checkInLogRequest)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			for _, invitee := range _invitees {
				if invitee.Name == name {
					checkInLogRequest.InviteeID = invitee.InviteeID
					_, err = db.CreateCheckInLog(context.Background(), checkInLogRequest)
					if err != nil {
						log.Fatalln(err)
					}
					break
				}
			}
		}
	}
}
