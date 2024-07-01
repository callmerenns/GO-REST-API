package entity

type Enrollment struct {
	UserID    uint    `gorm:"primaryKey;column:user_id"`
	ProductID uint    `gorm:"primaryKey;column:product_id"`
	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
