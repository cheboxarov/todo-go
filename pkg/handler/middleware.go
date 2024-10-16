package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "bad authorization", "header not found")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "bad authorization", "bad header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "bad authorization", err.Error())
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user not found")
	}
	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("bad user id")
	}
	return idInt, nil
}
