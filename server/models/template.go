package models

type Template struct {
	Id      int64  `gorm:"column:id"`
	Name    string `gorm:"column:name"`
	Content []byte `gorm:"column:content"`
}

func (Template) TableName() string {
	return "templates"
}
