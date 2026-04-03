package entity

type Category struct {
	ID   int32  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);unique;not null"`
	Slug string `gorm:"type:varchar(255);unique;not null;index:idx_categories_slug"`
}

func (Category) TableName() string {
	return "categories"
}
