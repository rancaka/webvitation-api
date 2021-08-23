package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/db"
	"github.com/rancaka/webvitation-web/models"
)

func GetEventsHandler(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	events, err := db.GetEventsByUID(c.Request().Context(), auth.UID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"List": events,
	})
}

func GetEventDetailHandler(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	slug := c.Param("eventSlug")

	event, err := db.GetEventBySlug(c.Request().Context(), slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if auth.UID != event.CreatorUID {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	event.Invitees, err = db.GetInviteesByEventID(c.Request().Context(), event.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event.Notifications, err = db.GetNotifications(c.Request().Context(), event.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Event": event,
	})
}

func CreateEventHandler(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	eventRequest := new(models.EventRequest)
	err := c.Bind(eventRequest)
	if err != nil {

		errors := []string{}
		if eventRequest.Slug == "" {
			errors = append(errors, "Slug is missing!")
		}

		if len(errors) > 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"Errors": errors,
			})
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eventRequest.CreatorUID = auth.UID
	newEventID, err := db.CreateEvent(c.Request().Context(), eventRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	duplicateInvitees, err := db.CreateInviteesBulk(c.Request().Context(), newEventID, eventRequest.Invitees)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event, err := db.GetEventByEventID(c.Request().Context(), newEventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event.Invitees, err = db.GetInviteesByEventID(c.Request().Context(), event.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Event":             event,
		"DuplicateInvitees": duplicateInvitees,
	})
}
