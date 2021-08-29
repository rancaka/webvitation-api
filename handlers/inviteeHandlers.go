package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-api/db"
	"github.com/rancaka/webvitation-api/models"
)

func CreateInviteesBulkHandler(c echo.Context) error {

	// auth := GetAuth(c)
	// if auth == nil {
	// 	return echo.NewHTTPError(http.StatusUnauthorized)
	// }

	slug := c.Param("eventSlug")

	event, err := db.GetEventBySlug(c.Request().Context(), slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	request := &struct {
		InviteeRequests []*models.InviteeRequest
	}{}

	err = c.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	duplicateInvitees, err := db.CreateInviteesBulk(c.Request().Context(), event.EventID, request.InviteeRequests)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"DuplicateInvitees": duplicateInvitees,
	})
}
