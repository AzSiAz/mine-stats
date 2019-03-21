package handler

import (
	"mine-stats/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (h *Handler) LoginHandler(c echo.Context) error {
	var form AuthForm
	err := c.Bind(&form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error":   err.Error(),
			"message": "Error on login",
		})
	}

	user, err := h.Store.VerifyLogin(form.Username, form.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error":   err.Error(),
			"message": "Error verifying who you are",
		})
	}

	user, err = h.Store.UpdateUserSessionIDAdd(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"message": "Error creating login session",
		})
	}

	sessionIDCookie := new(http.Cookie)
	sessionIDCookie.Value = user.SessionID
	sessionIDCookie.Name = "sessionID"
	if h.Prod {
		sessionIDCookie.Expires = time.Now().Add(24 * time.Hour * 31)
		sessionIDCookie.Secure = true
		sessionIDCookie.HttpOnly = true
	}

	c.SetCookie(sessionIDCookie)

	return c.JSON(http.StatusOK, Response{
		"username": user.Username,
		"id":       user.ID,
	})
}

func (h *Handler) LogoutHandler(c echo.Context) error {
	sessionIDCookie := new(http.Cookie)
	sessionIDCookie.Value = ""
	sessionIDCookie.Name = "sessionID"
	if h.Prod {
		sessionIDCookie.Expires = time.Now().Add(-1 * time.Hour)
		sessionIDCookie.Secure = true
		sessionIDCookie.HttpOnly = true
	}

	user := c.Get("user").(*models.User)
	_, err := h.Store.UpdateUserSessionIDRemove(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error":   err.Error(),
			"message": "Error verifying who you are",
		})
	}

	c.SetCookie(sessionIDCookie)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignUpHandler(c echo.Context) error {
	//c.String(http.StatusOK, "sign up")
	var form AuthForm
	err := c.Bind(&form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error":   err.Error(),
			"message": "Error signing up",
		})
	}

	_, err = h.Store.AddUser(form.Username, form.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"Error":   err.Error(),
			"message": "Error signing up with this info, try again",
		})
	}

	return c.JSON(http.StatusOK, Response{
		"error":   false,
		"message": "Sign Up successful try login now",
	})
}

func (h *Handler) MeHandler(c echo.Context) error {
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, Response{
		"ID":       user.ID,
		"username": user.Username,
	})
}
