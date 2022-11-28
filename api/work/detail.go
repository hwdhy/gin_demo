package work

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
	Id         int    `json:"id"`
	Type       int    `json:"type"`
	Title      string `json:"title"`
	Address    string `json:"address"`
	Edu        string `json:"edu"`
	Expre      string `json:"expre"`
	Num        string `json:"num"`
	AddTime    string `json:"add_time"`
	UpdateTime string `json:"update_time"`
	Sortindex  int    `json:"sortindex"`
	AdminId    int    `json:"admin_id"`
	Content    string `json:"content"`
	IsTese     int    `json:"is_tese"`
	Keyword    string `json:"keyword"`
	Desc       string `json:"desc"`
	Status     int    `json:"status"`
	Field1     string `json:"field_1" gorm:"column:field_1"`
	Field2     int    `json:"field_2" gorm:"column:field_2"`
}

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

	// 获取多语言类型
	language, _ := c.Get("language")
	languageInt := language.(int)

	// 未输入语言字段，默认为中文
	if languageInt == 0 {
		languageInt = 1
	}
	LogClient.WithField("language", languageInt).Infof("request_body: %+v", detailInput)

	// 数据库查询
	db := base.SlaveDb.Model(models.Work{})

	var workData models.Work
	if err := db.Where("id = ?", detailInput.ID).First(&workData).Error; err != nil && err != gorm.ErrRecordNotFound {
		LogClient.Errorf("find news list err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "find news list error",
		})
		return
	}

	resData := DetailOutputData{
		Id:        workData.Id,
		Type:      workData.Type,
		Title:     workData.Title,
		Address:   workData.Address,
		Edu:       workData.Edu,
		Expre:     workData.Expre,
		Num:       workData.Num,
		Sortindex: workData.Sortindex,
		AdminId:   workData.AdminId,
		Content:   workData.Content,
		IsTese:    workData.IsTese,
		Keyword:   workData.Keyword,
		Desc:      workData.Desc,
		Status:    workData.Status,
	}
	if languageInt == 1 {
		resData.UpdateTime = workData.UpdateTime.Format("06/01/02")
	} else {
		resData.UpdateTime = workData.UpdateTime.Format("01/02/06")
	}
	res := DetailOutput{
		Code: http.StatusOK,
		Data: resData,
	}
	LogClient.Infof("return: %+v", res)
	c.JSON(http.StatusOK, res)
}
