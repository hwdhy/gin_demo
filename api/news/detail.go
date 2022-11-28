package news

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"overseas-official-website/base"
	"overseas-official-website/models"
)

type DetailInput struct {
	ID int `json:"id" binding:"required"`
}

type DetailOutput struct {
	Code int              `json:"code"`
	Data DetailOutputData `json:"data"`
}

type DetailOutputData struct {
	ID         int    `json:"id"`
	Type       int    `json:"type"`
	Title      string `json:"title"`
	SubTitle   string `json:"sub_title"`
	Author     string `json:"author"`
	UpdateTime string `json:"update_time"`
	Sortindex  int    `json:"sortindex"`
	AdminId    int    `json:"admin_id"`
	Content    string `json:"content"`
	Click      int    `json:"click"`
	Comments   int    `json:"comments"`
	Image      string `json:"image"`
	Des        string `json:"des"`
	Status     int8   `json:"status"`
	Language   int8   `json:"language"`
	Rid        int    `json:"rid"`
	NewsType   int    `json:"news_type"`
}

// 新闻列表
func Detail(c *gin.Context) {
	// 获取日志对象
	LogClient := base.MultipleLog
	value, exists := c.Get("log")
	if exists {
		LogClient = value.(*logrus.Logger)
	}

	// 解析输入参数
	var detailInput DetailInput
	if err := c.Bind(&detailInput); err != nil {
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

	LogClient.WithField("language", languageInt).Infof("request_body: %+v", detailInput)

	var newsDetailData models.News
	// 数据库查询
	if err := base.SlaveDb.Model(models.News{}).Where("id = ?", detailInput.ID).First(&newsDetailData).Error; err != nil && err != gorm.ErrRecordNotFound {
		LogClient.Errorf("find news list err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "find news list error",
		})
		return
	}

	resData := DetailOutputData{
		ID:        newsDetailData.ID,
		Type:      newsDetailData.Type,
		Title:     newsDetailData.Title,
		SubTitle:  newsDetailData.SubTitle,
		Author:    newsDetailData.Author,
		Sortindex: newsDetailData.Sortindex,
		AdminId:   newsDetailData.AdminId,
		Content:   newsDetailData.Content,
		Click:     newsDetailData.Click,
		Comments:  newsDetailData.Comments,
		Image:     newsDetailData.Image,
		Des:       newsDetailData.Des,
		Status:    newsDetailData.Status,
		Language:  newsDetailData.Language,
		Rid:       newsDetailData.Rid,
		NewsType:  newsDetailData.NewsType,
	}
	if languageInt == 1 {
		resData.UpdateTime = newsDetailData.UpdateTime.Format("06/01/02")
	} else {
		resData.UpdateTime = newsDetailData.UpdateTime.Format("01/02/06")
	}
	res := DetailOutput{
		Code: http.StatusOK,
		Data: resData,
	}
	LogClient.Infof("return: %+v", res)
	c.JSON(http.StatusOK, res)
}
