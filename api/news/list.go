package news

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"overseas-official-website/base"
	"overseas-official-website/models"
)

type ListInput struct {
	PageNum  int    `json:"page_num" binding:"required"`
	PageSize int    `json:"page_size"`
	Type     int    `json:"type"`
	Keywords string `json:"keywords"`
}

type ListOutput struct {
	Code  int          `json:"code"`
	Count int64        `json:"count"`
	Data  []OutputData `json:"data"`
}

type OutputData struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	UpdateTime string `json:"update_time"`
	Image      string `json:"image"`
	Des        string `json:"des"`
}

// 新闻列表
func Lists(c *gin.Context) {
	// 获取日志对象
	LogClient := base.MultipleLog
	value, exists := c.Get("log")
	if exists {
		LogClient = value.(*logrus.Logger)
	}

	// 解析输入参数
	var listInput ListInput
	if err := c.Bind(&listInput); err != nil {
		LogClient.Errorf("parameter parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "parameter parsing error",
		})
		return
	}

	language, _ := c.Get("language")
	languageInt := language.(int)

	// 未输入语言字段，默认为中文
	if languageInt == 0 {
		languageInt = 1
	}

	LogClient.WithField("language", languageInt).Infof("request_body: %+v", listInput)

	// 数据库查询
	selectField := "id, title, update_time, image, des"
	db := base.SlaveDb.Model(models.News{}).Select(selectField).Where("language = ?", languageInt)
	if listInput.Type != 0 {
		db.Where("type = ?", listInput.Type)
	}
	if listInput.Keywords != "" {
		db.Where("title like ?", fmt.Sprintf("%%%s%%", listInput.Keywords))
	}

	if listInput.PageSize == 0 {
		listInput.PageSize = 10
	}

	offset := (listInput.PageNum - 1) * listInput.PageSize
	var count int64
	var newsData []models.News
	if err := db.Count(&count).Order("id desc").
		Offset(offset).Limit(listInput.PageSize).Find(&newsData).Error; err != nil && err != gorm.ErrRecordNotFound {

		LogClient.Errorf("find news list err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "find news list error",
		})
		return
	}

	resData := make([]OutputData, len(newsData))
	for k, news := range newsData {
		var updateTime string
		if languageInt == 1 {
			updateTime = news.UpdateTime.Format("06/01/02")
		} else {
			updateTime = news.UpdateTime.Format("01/02/06")
		}
		resData[k].ID = news.ID
		resData[k].Image = news.Image
		resData[k].UpdateTime = updateTime
		resData[k].Title = news.Title
		resData[k].Des = news.Des
	}

	res := ListOutput{
		Code:  http.StatusOK,
		Count: count,
		Data:  resData,
	}
	LogClient.Infof("return: %+v", res)
	c.JSON(http.StatusOK, res)
}
