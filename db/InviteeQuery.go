package db

import (
	"context"
	"log"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/rancaka/webvitation-api/models"
)

const inviteeColumns = `
	invitee_id,
	name,
	slug,
	event_id,
	created_at,
	updated_at
`

func getInviteesByEventID(ctx context.Context, eventID string) ([]*models.Invitee, error) {

	invitees := []*models.Invitee{}
	query := "SELECT " + inviteeColumns + " FROM event_invitees WHERE event_id = ?"
	err := sqlx.SelectContext(ctx, db, &invitees, query, eventID)
	if err != nil {
		log.Printf("[ERROR] db.GetInviteesByEventID -> sqlx.SelectContext: (%v)\n", err)
		return nil, err
	}

	return invitees, nil
}

func getInviteeByEventIDAndInviteeID(ctx context.Context, eventID, inviteeID string) (*models.Invitee, error) {

	invitee := new(models.Invitee)
	query := "SELECT " + inviteeColumns + " FROM event_invitees WHERE event_id = ? AND invitee_id = ? LIMIT 1"
	err := sqlx.GetContext(ctx, db, invitee, query, eventID, inviteeID)
	if err != nil {
		log.Printf("[ERROR] db.GetInviteeByEventIDAndSlug -> sqlx.GetContext: (%v)\n", err)
		return nil, err
	}

	return invitee, nil
}

func getInviteeByEventIDAndSlug(ctx context.Context, eventID, inviteeSlug string) (*models.Invitee, error) {

	invitee := new(models.Invitee)
	query := "SELECT " + inviteeColumns + " FROM invitees WHERE event_id = ? AND slug = ? LIMIT 1"
	err := sqlx.GetContext(ctx, db, invitee, query, eventID, inviteeSlug)
	if err != nil {
		log.Printf("[ERROR] db.GetInviteeByEventIDAndSlug -> sqlx.GetContext: (%v)\n", err)
		return nil, err
	}

	return invitee, nil
}

func createInviteesBulk(ctx context.Context, eventID string, inviteeRequests []*models.InviteeRequest) (duplicateInvitees []*models.InviteeRequest, err error) {

	// get current invitees
	var (
		query = `
		SELECT
			name
		FROM
			invitees
		WHERE
			event_id = ?
		`
		currentInvitees = []string{}
	)
	err = sqlx.SelectContext(ctx, db, &currentInvitees, query, eventID)
	if err != nil {
		log.Printf("[ERROR] db.CreateInviteesBulk -> sqlx.SelectContext: (%v)\n", err)
		return
	}

	// remove duplicates
	inviteeReqMap := make(map[string]*models.InviteeRequest)
	for _, inviteeRequest := range inviteeRequests {
		inviteeReqMap[inviteeRequest.Name] = inviteeRequest
	}

	insertQuery := "INSERT INTO event_invitees (name, slug, event_id) VALUES "
	args := []string{}
	for _, inviteeRequest := range inviteeReqMap {

		// check duplicate invitees
		isDuplicate := false

		// check duplicates against existing invitees
		for _, currentInvitee := range currentInvitees {
			if currentInvitee == inviteeRequest.Name {
				duplicateInvitees = append(duplicateInvitees, inviteeRequest)
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			continue
		}

		insertQuery += "("
		argsLength := len(args)
		endIndex := argsLength + 3
		for i := argsLength; i < endIndex; i++ {
			insertQuery += "?"
			if i != endIndex-1 {
				insertQuery += ", "
			}
		}
		insertQuery += "),"

		// setup data
		inviteeRequest.EventID = eventID
		inviteeRequest.Slug = slug.Make(inviteeRequest.Name)
		args = append(args, inviteeRequest.Name, inviteeRequest.Slug, inviteeRequest.EventID)
	}

	_, err = db.ExecContext(ctx, insertQuery, args)
	if err != nil {
		log.Printf("[ERROR] db.CreateInviteesBulk -> db.ExecContext: (%v)\n", err)
		return
	}

	return
}

func createInviteesBulk2(ctx context.Context, eventID string, inviteeRequests []*models.InviteeRequest) error {

	// var slugMap = make(map[string]int)

	tx := db.MustBeginTx(ctx, nil)
	for _, inviteeRequest := range inviteeRequests {

		// setup data
		inviteeRequest.EventID = eventID
		// inviteeRequest.Slug = slug.Make(inviteeRequest.Name)

		// if counter, exist := slugMap[inviteeRequest.Slug]; exist {
		// 	counter++
		// 	inviteeRequest.Slug += fmt.Sprintf("-%d", counter)
		// 	slugMap[inviteeRequest.Slug] = counter
		// } else {
		// 	slugMap[inviteeRequest.Slug] = 0
		// }

		_, err := tx.NamedExec(
			`INSERT INTO invitees (name, slug, event_id) VALUES (:name, :slug, :event_id)`,
			&inviteeRequest,
		)
		if err != nil {
			log.Printf("[ERROR] db.CreateInviteesBulk -> tx.NamedExec: (%v)\n", err)
			return err
		}
	}

	err := tx.Commit()
	if err != nil {
		log.Printf("[ERROR] db.CreateInviteesBulk -> tx.Commit: (%v)\n", err)
		return err
	}

	return nil
}
