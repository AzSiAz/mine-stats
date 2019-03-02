package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"mine-stats/models"
	"net/http"
	"time"
)

type AuthForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (h *Handler) LoginHandler(c echo.Context) error {
	var form AuthForm
	err := c.Bind(&form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
			"message": "Error on login",
		})
	}

	user, err := h.store.VerifyLogin(form.Username, form.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"message": "Error verifying who you are",
		})
	}

	user, err = h.store.UpdateUserWithSessionID(user)

	sessionIDCookie := new(http.Cookie)
	sessionIDCookie.Value = user.SessionID
	sessionIDCookie.Name = "sessionID"
	sessionIDCookie.Expires = time.Now().Add(24 * time.Hour * 31)
	sessionIDCookie.Secure = true
	sessionIDCookie.HttpOnly = true

	c.SetCookie(sessionIDCookie)

	return c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"id": user.ID,
	})
}

func (h *Handler) LogoutHandler(c echo.Context) error {
	sessionIDCookie := new(http.Cookie)
	sessionIDCookie.Value = ""
	sessionIDCookie.Name = "sessionID"
	sessionIDCookie.Expires = time.Now().Add(-1 * time.Hour)
	sessionIDCookie.Secure = true
	sessionIDCookie.HttpOnly = true

	c.SetCookie(sessionIDCookie)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignUpHandler(c echo.Context) error {
	//c.String(http.StatusOK, "sign up")
	var form AuthForm
	err := c.Bind(&form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
			"message": "Error signing up",
		})
	}

	_, err = h.store.AddUser(form.Username, form.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
			"message": "Error signing up with this info, try again",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"error": false,
		"message": "Sign Up successful try login now",
	})
}

func (h *Handler) MeHandler(c echo.Context) error {
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, gin.H{
		"ID": user.ID,
		"username": user.Username,
	})
}