package dto

import (
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
)

type EventReq struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Date        utils.DateOnly `json:"date" binding:"required"`
}

type EventRes struct {
	model.Event
}
