package utils

import (
	"car-park/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	DEFAULT_LIMIT = 5
	DEFAULT_PAGE  = 1
	DEFAUL_SORT   = "created_at asc"
)

func GenPaginationFromRequest(ctx *gin.Context) models.Pagination {
	limit := DEFAULT_LIMIT
	page := DEFAULT_PAGE
	sort := DEFAUL_SORT

	query := ctx.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break
		}

	}

	pagination := models.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

	return pagination
}
