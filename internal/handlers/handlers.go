package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reversersed/zyfra/internal/handlers/models"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Log in user with login and password
// @Tags         sessions
// @Produce      json
// @Param        body body models.LoginCommand true "User credentials"
// @Success      200  {object}   handlers.HandleLoginCommand.LoginResponse "Generated session key"
// @Failure      400  {object}  models.RequestError "Received bad request body"
// @Failure      404  {object}  models.RequestError "User was not found or password is incorrect"
// @Router       /sessions [post]
func (h *handler) HandleLoginCommand(c *gin.Context) {
	request := models.LoginCommand{}
	_ = c.BindJSON(&request)

	if len(request.Login) == 0 {
		c.JSON(http.StatusBadRequest, models.RequestError{Message: "Excepted non-empty login", Error: "login length was zero"})
		return
	}
	if len(request.Password) == 0 {
		c.JSON(http.StatusBadRequest, models.RequestError{Message: "Excepted non-empty password", Error: "password length was zero"})
		return
	}
	password, exist := h.cfg[request.Login]
	if !exist {
		c.JSON(http.StatusNotFound, models.RequestError{Message: "User does not exist", Error: "user not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(password, []byte(request.Password)); err != nil {
		c.JSON(http.StatusNotFound, models.RequestError{Message: "Received wrong password", Error: err.Error()})
		return
	}
	type LoginResponse struct {
		Session string `json:"session"`
	}
	c.JSON(http.StatusOK, &LoginResponse{Session: h.service.CreateSession()})
}

// @Summary      Authenticate user with session key
// @Tags         sessions
// @Produce      json
// @Param        session path string true "Session key"
// @Success      204 "Session is valid"
// @Failure      400  {object}  models.RequestError "Received bad request body"
// @Failure      401  {object}  models.RequestError "Session is invalid"
// @Router       /sessions/{session} [get]
func (h *handler) HandleAuthRequest(c *gin.Context) {
	request := models.AuthRequest{}
	_ = c.BindUri(&request)

	if err := h.service.CheckSession(request.Session); err != nil {
		c.JSON(http.StatusUnauthorized, models.RequestError{Message: "User not authorized", Error: "Session not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      Delete session by provided key
// @Tags         sessions
// @Produce      json
// @Param        session path string true "Session key"
// @Success      204 "Session deleted"
// @Failure      400  {object}  models.RequestError "Received bad request body"
// @Failure      404  {object}  models.RequestError "Session not found"
// @Router       /sessions/{session} [delete]
func (h *handler) HandleDeleteCommand(c *gin.Context) {
	request := models.DeleteCommand{}
	_ = c.BindUri(&request)

	if err := h.service.Delete(request.Session); err != nil {
		c.JSON(http.StatusNotFound, models.RequestError{Message: "Session was not deleted", Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
