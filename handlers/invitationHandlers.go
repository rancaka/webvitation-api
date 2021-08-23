package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/db"
)

func GetInvitationPageHandler(c echo.Context) error {

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	media, err := db.GetMediaByEventID(c.Request().Context(), event.EventID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = c.Render(http.StatusOK, ".html", echo.Map{
		"Event": event,
		"Media": media,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func GetInvitationPageWithInviteeHandler(c echo.Context) error {

	eventSlug := c.Param("eventSlug")
	event, err := db.GetEventBySlug(c.Request().Context(), eventSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	inviteeSlug := c.Param("inviteeSlug")
	invitee, err := db.GetInviteeByEventIDAndSlug(c.Request().Context(), event.EventID, inviteeSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	photos, err := db.GetMediaByEventID(c.Request().Context(), event.EventID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = c.Render(http.StatusOK, ".html", echo.Map{
		"Event":   event,
		"Invitee": invitee,
		"Photos":  photos,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
