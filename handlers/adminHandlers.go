package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/db"
)

func GetAdminEventListPage(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	events, err := db.GetEventsByUID(c.Request().Context(), auth.UID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, "event_list.html", echo.Map{
		"Events": events,
	})
}

func GetAdminCreateEventPage(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	return Render(c, http.StatusOK, "create_event.html", nil)
}
