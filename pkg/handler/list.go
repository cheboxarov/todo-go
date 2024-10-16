package handler

import (
	"net/http"
	"strconv"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	var input todogo.TodoList
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Count int               `json:"count"`
	Data  []todogo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	c.JSON(200, getAllListsResponse{Count: len(lists), Data: lists})
}

type getByIdListResponse struct {
	Data todogo.TodoList `json:"data"`
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	c.JSON(200, getByIdListResponse{Data: list})
}

type statusResponse struct {
	Status string `json:"status"`
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	var input todogo.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	err = h.services.TodoList.Update(userId, id, input)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}
	c.JSON(200, statusResponse{Status: "Updated"})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something is wrong", err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something is wrong", err.Error())
		return
	}

	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		NewErrorResponse(c, 500, "something is wrong", err.Error())
		return
	}

	c.JSON(200, statusResponse{Status: "Deleted"})
}
