package main

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-web/handlers"
)

type MyTemplate struct {
	templates *template.Template
}

func (t *MyTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"Add": func(a, b int) int {
			return a + b
		},
		"EscapeHTML": func(json string) template.HTML {
			return template.HTML(json)
		},
		"GetDay": func(date time.Time) string {
			return date.Format("02")
		},
		"GetHourKitchen": func(date time.Time) string {
			return date.Format(time.Kitchen)
		},
		"GetHourAndMinute": func(date time.Time) string {

			if date.Format("2006-02-01 15:04:05") == "0000-01-01 00:00:00" {
				return ""
			}

			return date.Format("15:04")
		},
		"GetHour": func(date time.Time) int {
			return date.Hour()
		},
		"GetMonth": func(date time.Time) string {
			return date.Format("Jan")
		},
		"GetMonthNumber": func(date time.Time) int {
			return int(date.Month())
		},
		"GetWeekday": func(date time.Time) string {
			return date.Weekday().String()
		},
		"GetYear": func(date time.Time) int {
			return date.Year()
		},
	}).ParseGlob("views/**/*.html")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Renderer = &MyTemplate{
		templates: tmpl,
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			myContext := &handlers.MyContext{c}
			return next(myContext)
		}
	})

	e.Static("/static", "assets")

	e.GET("/login", handlers.GetLoginPageHandler)
	e.POST("/login", handlers.LoginHandler)

	admin := e.Group("/admin")
	admin.Use(handlers.VerifyToken)
	admin.GET("", handlers.GetAdminEventListPage)
	admin.GET("/event/create", handlers.GetAdminCreateEventPage)

	api := e.Group("/api")
	events := api.Group("/events")
	events.Use(handlers.VerifyToken)
	events.GET("", handlers.GetEventsHandler)
	events.POST("/create", handlers.CreateEventHandler)

	eventDetail := events.Group("/:eventSlug")
	eventDetail.GET("", handlers.GetEventDetailHandler)
	eventDetail.POST("/create-invitees", handlers.CreateInviteesBulkHandler)
	eventDetail.POST("/check-in", handlers.CheckInHandler)
	eventDetail.POST("/sync", handlers.SyncHandler)

	eventPage := e.Group("/:eventSlug")
	eventPage.GET("", handlers.GetInvitationPageHandler)
	eventPage.GET("/check-in", handlers.GetCheckInPageHandler, handlers.VerifyToken)
	eventPage.GET("/check-in-list", handlers.GetCheckInListPageHandler, handlers.VerifyToken)
	eventPage.GET("/to/:inviteeSlug", handlers.GetInvitationPageWithInviteeHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
