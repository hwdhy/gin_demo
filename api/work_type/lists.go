package work_type

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"overseas-official-website/base"
	"overseas-official-website/models"
)

// 输入参数
type ListInput struct {
	PageNum  int    `json:"page_num" binding:"required"`
	PageSize int    `json:"page_size"`
	Keywords string `json:"keywords"`
}

type ListOutput struct {
	Code  int          `json:"code"`
	Count int64        `json:"count"`
	Data  []OutputData `json:"data"`
}

type OutputData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
	selectField := "id, name, en_name"
	db := base.SlaveDb.Model(models.WorkType{}).Select(selectField)
	if listInput.Keywords != "" {
		db.Where("name like ?", fmt.Sprintf("%%%s%%", listInput.Keywords))
	}

	if listInput.PageSize == 0 {
		listInput.PageSize = 10
	}

	offset := (listInput.PageNum - 1) * listInput.PageSize
	var Count int64
	var workTypeData []models.WorkType
	if err := db.Count(&Count).Order("id desc").Offset(offset).Limit(listInput.PageSize).Find(&workTypeData).Error; err != nil && err != gorm.ErrRecordNotFound {
		LogClient.Errorf("find news list err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "find news list error",
		})
		return
	}

	resData := make([]OutputData, len(workTypeData))
	for k, workType := range workTypeData {
		resData[k].ID = workType.ID
		// 中英文判断
		if languageInt == 1 {
			resData[k].Name = workType.Name
		} else {
			resData[k].Name = workType.EnName
		}
	}

	res := ListOutput{
		Code:  http.StatusOK,
		Count: Count,
		Data:  resData,
	}
	LogClient.Infof("return: %+v", res)
	c.JSON(http.StatusOK, res)
}
