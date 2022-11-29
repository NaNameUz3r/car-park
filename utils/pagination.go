package utils

import (
	"car-park/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	DEFAULT_PAGINATION_LIMIT = 10
	DEFAULT_PAGINATION_PAGE  = 1
	DEFAUL_PAGINATION_SORT   = "id desc"
)

func GenPaginationFromRequest(ctx *gin.Context) models.Pagination {
	limit := DEFAULT_PAGINATION_LIMIT
	page := DEFAULT_PAGINATION_PAGE
	sort := DEFAUL_PAGINATION_SORT

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
