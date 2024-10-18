package handler

import (
	"net/http"
	"strconv"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, 500, "invalid list id param", "")
		return
	}

	var input todogo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		return
	}
	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllItemsResponse struct {
	Count int               `json:"count"`
	Data  []todogo.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllItemsResponse{
		Count: len(items),
		Data:  items,
	})
}

type getItemByIdResponse struct {
	Data todogo.TodoItem `json:"data"`
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	item, err := h.services.TodoItem.GetById(userId, listId, itemId)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	c.JSON(200, getItemByIdResponse{
		Data: item,
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	c.JSON(200, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	var input todogo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error(), err.Error())
		return
	}
	if err := input.Validate(); err != nil {
		NewErrorResponse(c, 400, err.Error(), err.Error())
		return
	}

	if err := h.services.TodoItem.Update(userId, itemId, input); err != nil {
		NewErrorResponse(c, 500, "", err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
