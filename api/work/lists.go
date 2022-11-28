package work

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
	Id         int    `json:"id"`
	Type       int    `json:"type"`
	Title      string `json:"title"`
	Address    string `json:"address"`
	UpdateTime string `json:"update_time"`
}

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

	// 获取多语言类型
	language, _ := c.Get("language")
	languageInt := language.(int)

	// 未输入语言字段，默认为中文
	if languageInt == 0 {
		languageInt = 1
	}
	LogClient.WithField("language", languageInt).Infof("request_body: %+v", listInput)

	// 数据库查询
	selectField := "id, type, title, address, update_time"
	db := base.SlaveDb.Model(models.Work{}).Select(selectField)
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
	var workData []models.Work
	if err := db.Count(&count).Order("id desc").Offset(offset).Limit(listInput.PageSize).Find(&workData).Error; err != nil && err != gorm.ErrRecordNotFound {
		LogClient.Errorf("find news list err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "find news list error",
		})
		return
	}

	resData := make([]OutputData, len(workData))
	for k, work := range workData {
		var updateTime string
		if languageInt == 1 {
			updateTime = work.UpdateTime.Format("06/01/02")
		} else {
			updateTime = work.UpdateTime.Format("01/02/06")
		}

		resData[k].Id = work.Id
		resData[k].Type = work.Type
		resData[k].Title = work.Title
		resData[k].Address = work.Address
		resData[k].UpdateTime = updateTime
	}

	res := ListOutput{
		Code:  http.StatusOK,
		Count: count,
		Data:  resData,
	}
	LogClient.Infof("return: %+v", res)
	c.JSON(http.StatusOK, res)
}
