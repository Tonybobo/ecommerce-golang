package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Pagination(c *gin.Context) (page int, limit int) {
	var (
		pageNum  = 1
		limitNum = 30
	)

	if c.Query("page") != "" {
		pageNum, _ = strconv.Atoi(c.Query("page"))
	}
	if c.Query("limit") != "" {
		limitNum, _ = strconv.Atoi(c.Query("limit"))
	}

	return pageNum, limitNum
}
