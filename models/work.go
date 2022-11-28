package models

import "time"

type Work struct {
	Id         int       `json:"id"`
	Type       int       `json:"type"`
	Title      string    `json:"title"`
	Address    string    `json:"address"`
	Edu        string    `json:"edu"`
	Expre      string    `json:"expre"`
	Num        string    `json:"num"`
	AddTime    time.Time `json:"add_time"`
	UpdateTime time.Time `json:"update_time"`
	Sortindex  int       `json:"sortindex"`
	AdminId    int       `json:"admin_id"`
	Content    string    `json:"content"`
	IsTese     int       `json:"is_tese"`
	Keyword    string    `json:"keyword"`
	Desc       string    `json:"desc"`
	Status     int       `json:"status"`
	Field1     string    `json:"field_1" gorm:"column:field_1"`
	Field2     int       `json:"field_2" gorm:"column:field_2"`
}

func (*Work) TableName() string {
	return "work"
}
