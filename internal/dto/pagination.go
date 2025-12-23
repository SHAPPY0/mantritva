package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/pkg/utils"
)

type Pagination struct {
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

func NewPagination(c *gin.Context) Pagination {
	pagination := Pagination{}
	var limit string
	var offset string

	if offset = c.Query("offset"); offset == "" {
		offset = "0"
	}
	if limit = c.Query("limit"); limit == "" {
		limit = "20"
	}
	limitInt, err := utils.StrToInt(limit)

	if err != nil {
		utils.Error(c, 400, "Invalid limit", err)
		return pagination
	}
	offsetInt, err := utils.StrToInt(offset)
	if err != nil {
		utils.Error(c, 400, "Invalid offset", err)
		return pagination
	}
	pagination.Limit = limitInt
	pagination.Offset = offsetInt
	return pagination
}
