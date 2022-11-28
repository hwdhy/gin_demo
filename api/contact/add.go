package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"net/http"
	"overseas-official-website/base"
	"overseas-official-website/models"
	"time"
)

type AddInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Msg   string `json:"msg"`
}

type AddOutput struct {
	Code int `json:"code"`
}

func Add(c *gin.Context) {
	// 获取日志对象
	LogClient := base.MultipleLog
	value, exists := c.Get("log")
	if exists {
		LogClient = value.(*logrus.Logger)
	}

	// 解析输入参数
	var addInput AddInput
	if err := c.Bind(&addInput); err != nil {
		LogClient.Errorf("parameter parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "parameter parsing error",
		})
		return
	}
	LogClient.Infof("request_body: %+v", addInput)

	// 新增数据
	contactData := models.Contact{
		Name:    addInput.Name,
		Email:   addInput.Email,
		Msg:     html.EscapeString(addInput.Msg),
		IP:      c.ClientIP(),
		AddTime: time.Now(),
	}

	if err := base.MasterDb.Model(models.Contact{}).Create(&contactData).Error; err != nil {
		LogClient.Errorf("add contact data err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "server error",
		})
		return
	}

	resData := &AddOutput{
		Code: http.StatusOK,
	}
	c.JSON(http.StatusOK, resData)
}
