package crudex

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: Complete this part of the code to allow search and pagination for the models
const (
	SEARCH_OP_EQUAL           = "="
	SEARCH_OP_NOT_EQUAL       = "<>"
	SEARCH_OP_MATCHES_PATTERN = "~"
	SEARCH_OP_IN              = "in"
	SEARCH_OP_NOT_IN          = "not in"
	SEARCH_OP_LT              = "<"
	SEARCH_OP_LTE             = "<="
	SEARCH_OP_GT              = ">"
	SEARCH_OP_GTE             = ">="
)

// Struct for filter and pagination
type SearchArgs struct {
	Search string
	Page   int
	Limit  int
}

func NewSearchArgs() SearchArgs {
	return SearchArgs{
		Search: "",
		Page:   1,
		Limit:  500,
	}
}

func NewSearchArgsFromQuery(c *gin.Context) (SearchArgs, error) {
	filter := NewSearchArgs()
	filter.Search = c.Query("search")

	if pageStr, exists := c.GetQuery("page"); exists {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			filter.Page = page
		} else {
			return filter, err
		}
	}
	if limitStr, exists := c.GetQuery("limit"); exists {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			filter.Limit = limit
		} else {
			return filter, err
		}
	}
	return filter, nil
}
