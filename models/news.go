package models

import "time"

type News struct {
	ID         int       `json:"id"`
	Type       int       `json:"type"`
	Title      string    `json:"title"`
	SubTitle   string    `json:"sub_title"`
	Author     string    `json:"author"`
	AddTime    time.Time `json:"add_time"`
	UpdateTime time.Time `json:"update_time"`
	Sortindex  int       `json:"sortindex"`
	AdminId    int       `json:"admin_id"`
	Content    string    `json:"content"`
	Click      int       `json:"click"`
	Comments   int       `json:"comments"`
	Image      string    `json:"image"`
	Des        string    `json:"des"`
	Status     int8      `json:"status"`
	Language   int8      `json:"language"`
	Rid        int       `json:"rid"`
	NewsType   int       `json:"news_type"`
}

func (*News) TableName() string {
	return "news"
}
