package dto

import (
	"database/sql"
	"time"

	"github.com/altsaqif/go-rest/cmd/entity"
	"gorm.io/gorm"
)

type AuthRequestRegisterDto struct {
	FirstName       string `gorm:"type:varchar(300);not null" json:"firstname"`
	LastName        string `gorm:"type:varchar(300);not null" json:"lastname"`
	Email           string `gorm:"not null;unique" json:"email"`
	Password        string `gorm:"not null" json:"password"`
	PasswordConfirm string `gorm:"not null" json:"password_confirm"`
	Role            string `gorm:"not null" json:"role"`
}

type AuthRequestLoginDto struct {
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type AuthResponseDto struct {
	Token string `json:"token"`
}

type AuthResponseRegisterDto struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(300);not null" json:"firstname"`
	LastName  string `gorm:"type:varchar(300);not null" json:"lastname"`
	Email     string `gorm:"not null;unique" json:"email"`
	Role      string `gorm:"not null" json:"role"`
}

type DeletedAt sql.NullTime

type UserWithoutProducts struct {
	ID        uint        `json:"ID"`
	CreatedAt time.Time   `json:"CreatedAt"`
	UpdatedAt time.Time   `json:"UpdatedAt"`
	DeletedAt DeletedAt   `gorm:"index" json:"DeletedAt,omitempty"`
	FirstName string      `json:"firstname"`
	LastName  string      `json:"lastname"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Role      string      `json:"role"`
	Products  interface{} `json:"products,omitempty"`
}

type ProductWithUsers struct {
	ID          uint                  `json:"ID"`
	CreatedAt   time.Time             `json:"CreatedAt"`
	UpdatedAt   time.Time             `json:"UpdatedAt"`
	DeletedAt   DeletedAt             `gorm:"index" json:"DeletedAt,omitempty"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Stock       int                   `json:"stock"`
	Price       float64               `json:"price"`
	Users       []UserWithoutProducts `json:"users"`
}

type ProductWithoutUsers struct {
	ID          uint      `json:"ID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	DeletedAt   DeletedAt `json:"DeletedAt,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	Price       float64   `json:"price"`
}

type UserWithProducts struct {
	ID        uint                  `json:"ID"`
	CreatedAt time.Time             `json:"CreatedAt"`
	UpdatedAt time.Time             `json:"UpdatedAt"`
	DeletedAt DeletedAt             `json:"DeletedAt,omitempty"`
	FirstName string                `json:"firstname"`
	LastName  string                `json:"lastname"`
	Email     string                `json:"email"`
	Password  string                `json:"password"`
	Role      string                `json:"role"`
	Products  []ProductWithoutUsers `json:"products"`
}

// Helper function to convert Product model to ProductWithUsers DTO
func ConvertProductToResponse(product entity.Product) ProductWithUsers {
	responseProduct := ProductWithUsers{
		ID:          product.ID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		DeletedAt:   DeletedAt(product.DeletedAt),
		Name:        product.Name,
		Description: product.Description,
		Stock:       product.Stock,
		Price:       product.Price,
		Users:       []UserWithoutProducts{},
	}

	for _, user := range product.Users {
		responseProduct.Users = append(responseProduct.Users, UserWithoutProducts{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: DeletedAt(user.DeletedAt),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			Role:      user.Role,
		})
	}

	return responseProduct
}

// Helper function to convert User model to UserWithProducts DTO
func ConvertUserToResponse(user entity.User) UserWithProducts {
	responseUser := UserWithProducts{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: DeletedAt(user.DeletedAt),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Products:  []ProductWithoutUsers{},
	}

	for _, product := range user.Products {
		responseUser.Products = append(responseUser.Products, ProductWithoutUsers{
			ID:          product.ID,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			DeletedAt:   DeletedAt(product.DeletedAt),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
		})
	}

	return responseUser
}
