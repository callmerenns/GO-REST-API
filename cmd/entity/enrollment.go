package entity

type Enrollment struct {
	UserID    uint `gorm:"primaryKey;column:user_id"`
	ProductID uint `gorm:"primaryKey;column:product_id"`
}
