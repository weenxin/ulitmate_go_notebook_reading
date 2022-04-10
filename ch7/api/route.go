package api

import "github.com/gin-gonic/gin"

func InitRoute(group *gin.RouterGroup) {
	group.GET("/:book_id", getBook)
	group.POST("/", CreateBook)
}
