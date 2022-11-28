package models

import "time"

// WorkType 职位类型
type WorkType struct {
	ID         int       `json:"id"`
	EnName     string    `json:"en_name"`
	Name       string    `json:"name"`
	Sortindex  int       `json:"sortindex"`
	Tags       string    `json:"tags"`
	AddTime    time.Time `json:"add_time"`
	UpdateTime time.Time `json:"update_time"`
	AdminID    int       `json:"admin_id"`
}

func (*WorkType) TableName() string {
	return "work_type"
}
