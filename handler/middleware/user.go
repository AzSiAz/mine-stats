package middleware

import (
	"github.com/labstack/echo/v4"
	"mine-stats/models"
	"mine-stats/store"
	"net/http"
)

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	st := store.GetStore()
	return func(c echo.Context) error {
		cookieSessionID, err := c.Cookie("sessionID")
		if err != nil {
			return c.NoContent(http.StatusForbidden)
		}

		user, err := st.GetUserBySessionID(cookieSessionID.Value)
		if err != nil {
			return c.NoContent(http.StatusForbidden)
		}

		c.Set("user", user)
		return next(c)
	}
}

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*models.User)

		if user.Role != models.AdminRole {
			return c.NoContent(http.StatusForbidden)
		}

		return next(c)
	}
}
