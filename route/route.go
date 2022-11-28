package route

import (
	"github.com/gin-gonic/gin"
	"overseas-official-website/api/contact"
	"overseas-official-website/api/news"
	"overseas-official-website/api/work"
	"overseas-official-website/api/work_type"
)

func InitRouters(engine *gin.Engine) {
	newsGroup := engine.Group("/news")
	newsGroup.POST("/list", news.Lists)
	newsGroup.POST("/detail", news.Detail)

	workGroup := engine.Group("/work")
	workGroup.POST("/list", work.Lists)
	workGroup.POST("/detail", work.Detail)

	workTypeGroup := engine.Group("/work_type")
	workTypeGroup.POST("/list", work_type.Lists)

	contactGroup := engine.Group("/contact")
	contactGroup.POST("/add", contact.Add)
}
