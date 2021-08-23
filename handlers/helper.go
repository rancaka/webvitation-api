package handlers

import (
	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/models"
)

type MyContext struct {
	echo.Context
}

func (c *MyContext) GetHeader(key string) string {
	return c.Request().Header.Get(key)
}

func GetAuth(c echo.Context) *models.Auth {

	auth, ok := c.Get("AUTH").(*models.Auth)
	if !ok {
		return nil
	}

	return auth
}

func Render(c echo.Context, code int, name string, data interface{}) error {

	auth := GetAuth(c)

	return c.Render(code, name, echo.Map{
		"Auth": auth,
		"Data": data,
	})
}
