package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/db"
	"github.com/rancaka/webvitation-web/models"
)

func GetCheckInPageHandler(c echo.Context) error {

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	invitees, err := db.GetInviteesByEventID(c.Request().Context(), event.EventID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = c.Render(http.StatusOK, "checkin_basic.html", echo.Map{
		"Event":    event,
		"Invitees": invitees,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func GetCheckInListPageHandler(c echo.Context) error {

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = c.Render(http.StatusOK, "checkin_list_basic.html", echo.Map{
		"Event": event,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func CheckInHandler(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	checkIn := new(models.CheckIn)
	err := c.Bind(checkIn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	checkIn.CreatorUID = auth.UID

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if checkIn.InviteeName == "" {
		invitee, err := db.GetInviteeByEventIDAndInviteeID(c.Request().Context(), event.EventID, checkIn.InviteeID)
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		checkIn.InviteeName = invitee.Name
	}

	notification := &models.Notification{
		EventID: event.EventID,
		Type:    models.NotificationCheckIn,
		CheckIn: checkIn,
	}

	notification.NotificationID, err = db.CreateNotification(c.Request().Context(), notification)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// checkInLog := models.Notification{
	// 	CheckInLogID: newCheckInLogID,
	// 	checkIn:      *checkIn,
	// }

	// go func() {
	// 	insertToFirebase(checkInLog)
	// }()

	return c.JSON(http.StatusOK, notification)
}

func SyncHandler(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	type SyncReq struct {
		CheckIns []*models.CheckIn
	}

	type SyncRes struct {
		Notifications []*models.Notification
	}

	syncReq := new(SyncReq)
	err := c.Bind(syncReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	syncRes := new(SyncRes)
	for _, checkIn := range syncReq.CheckIns {

		checkIn.CreatorUID = auth.UID
		if checkIn.InviteeName == "" {
			invitee, err := db.GetInviteeByEventIDAndInviteeID(c.Request().Context(), event.EventID, checkIn.InviteeID)
			if err != nil && err != sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			checkIn.InviteeName = invitee.Name
		}

		notificaton := &models.Notification{
			EventID: event.EventID,
			Type:    models.NotificationCheckIn,
		}

		notificaton.NotificationID, err = db.CreateNotification(c.Request().Context(), notificaton)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		syncRes.Notifications = append(syncRes.Notifications, notificaton)

		// go func(checkInLog models.CheckInLog) {
		// 	insertToFirebase(checkInLog)
		// }(checkInLog)
	}

	return c.JSON(http.StatusOK, syncRes)
}

// func insertToFirebase(checkInLog models.CheckInLog) {
// 	tempTime, err := time.Parse(db.CheckInAtLayout, checkInLog.CheckedInAt)
// 	if err != nil {
// 		fmt.Printf("[ERROR] handlers.CheckInHandler ->  parsing time: (%v)\n", err)
// 	}

// 	checkInLog.CheckedInAtUnix = tempTime.Unix()
// 	path := fmt.Sprintf("events/%s/check-ins", checkInLog.EventID)
// 	_, err = fb.Firestore.Collection(path).Doc(checkInLog.CheckInLogID).Set(context.Background(), checkInLog)
// 	if err != nil {
// 		fmt.Printf("[ERROR] handlers.CheckInHandler ->  Set data to firestore: (%v)\n", err)
// 	}
// }
