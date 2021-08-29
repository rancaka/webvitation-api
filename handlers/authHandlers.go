package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-api/models"
)

func GetLoginPageHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func LoginHandler(c echo.Context) error {

	authRequest := new(struct {
		FirebaseAuthToken string
	})

	err := c.Bind(authRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	decoded, err := fb.Auth.VerifyIDToken(c.Request().Context(), authRequest.FirebaseAuthToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Return error if the sign-in is older than 5 minutes.
	authTime := decoded.Claims["auth_time"].(float64)
	now := float64(time.Now().Unix())
	if now-authTime > 5*60 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Login is expired")
	}

	expiresIn := time.Hour * 24 * 7
	cookie, err := fb.Auth.SessionCookie(c.Request().Context(), authRequest.FirebaseAuthToken, expiresIn)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:   "session",
		Value:  cookie,
		MaxAge: int(expiresIn.Seconds()),
	})

	return c.JSON(http.StatusOK, authRequest)
}

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !config.IsProduction() {
			// bypass me
			auth := &models.Auth{
				UID:     "hIQ8ugRsjKXB6EApQGCUOD0rbJ83",
				Email:   "rancaka.dev@gmail.com",
				Name:    "Adityo Rancaka",
				Picture: "",
			}

			c.Set("AUTH", auth)
			return next(c)
		}

		var decodedToken *auth.Token
		cookie, err := c.Cookie("session")
		if err != nil {

			// check session from Authorization Header
			mc := c.(*MyContext)

			authorization := mc.GetHeader("Authorization")
			idToken := strings.TrimPrefix(authorization, "Bearer ")
			if idToken == "" {
				return RedirectToLoginPage(c)
			}

			decodedToken, err = fb.Auth.VerifyIDTokenAndCheckRevoked(mc.Request().Context(), idToken)
			if err != nil {
				return RedirectToLoginPage(c)
			}

		} else {

			decodedToken, err = fb.Auth.VerifySessionCookieAndCheckRevoked(c.Request().Context(), cookie.Value)
			if err != nil {
				return RedirectToLoginPage(c)
			}
		}

		email, ok := decodedToken.Claims["email"].(string)
		if !ok || email == "" {
			return RedirectToLoginPage(c)
		}

		name, ok := decodedToken.Claims["name"].(string)
		if !ok || name == "" {
			name = decodedToken.UID
		}

		picture, ok := decodedToken.Claims["picture"].(string)
		if !ok || picture == "" {
			picture = "https://t4.ftcdn.net/jpg/00/64/67/63/360_F_64676383_LdbmhiNM6Ypzb3FM4PPuFP9rHe7ri8Ju.jpg"
		}

		auth := &models.Auth{
			UID:     decodedToken.UID,
			Email:   email,
			Name:    name,
			Picture: picture,
		}

		c.Set("AUTH", auth)
		return next(c)
	}
}

func RedirectToLoginPage(c echo.Context) error {
	return c.Redirect(http.StatusFound, fmt.Sprintf("/login?redirectURL=%s", c.Request().RequestURI))
}
