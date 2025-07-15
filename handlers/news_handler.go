package handlers

import (
	"news-fetcher/api"

	"github.com/gin-gonic/gin"
)

func FetchNewsHandler(c *gin.Context) {
	api.FetchNews(c)
}
