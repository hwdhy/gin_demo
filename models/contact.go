package models

import "time"

// 联系我们
type Contact struct {
	ID      int       `json:"id" gorm:"primaryKey"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Msg     string    `json:"msg"`
	IP      string    `json:"ip"`
	AddTime time.Time `json:"add_time"`
}

func (*Contact) TableName() string {
	return "contact"
}
