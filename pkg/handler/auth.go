package handler

import (
	"net/http"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todogo.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "bad username or password", err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
