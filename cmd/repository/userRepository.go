package repository

import (
	"log"
	"math"

	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(payload entity.User) (dto.UserWithProducts, error)
	FindByID(id uint) (dto.UserWithProducts, error)
	FindByEmail(email string) (dto.UserWithProducts, error)
	FindAll(page, size int) ([]dto.UserWithProducts, model.Paging, error)
}

type userRepository struct {
	db *gorm.DB
}

// FindAll implements UserRepository.
func (u *userRepository) FindAll(page, size int) ([]dto.UserWithProducts, model.Paging, error) {
	var users []entity.User
	offset := (page - 1) * size

	type result struct {
		totalUsers int64
		err        error
	}

	resultChan := make(chan result)

	go func() {
		// Retrieve total count of users
		var totalUsers int64
		err := u.db.Model(&entity.User{}).Count(&totalUsers).Error
		resultChan <- result{totalUsers, err}
	}()

	var totalUsers int64
	res := <-resultChan
	if res.err != nil {
		log.Printf("userRepository.FindAll: Error counting users: %v \n", res.err)
		return nil, model.Paging{}, res.err
	}
	totalUsers = res.totalUsers

	// Retrieve paginated users
	if err := u.db.Limit(size).Offset(offset).Preload("Products").Find(&users).Error; err != nil {
		log.Printf("userRepository.FindAll: Error fetching users: %v \n", err)
		return nil, model.Paging{}, err
	}

	responseUsers := make([]dto.UserWithProducts, len(users))
	for i, user := range users {
		responseUsers[i] = dto.ConvertUserToResponse(user)
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   int(totalUsers),
		TotalPages:  int(math.Ceil(float64(totalUsers) / float64(size))),
	}

	return responseUsers, paging, nil
}

func (u *userRepository) Create(payload entity.User) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)

	go func() {
		if err := u.db.Create(&payload).Error; err != nil {
			resultChan <- result{dto.UserWithProducts{}, err}
			return
		}
		resultChan <- result{dto.ConvertUserToResponse(payload), nil}
	}()

	res := <-resultChan
	return res.user, res.err
}

func (u *userRepository) FindByID(id uint) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)

	go func() {
		var user entity.User
		if err := u.db.Preload("Products").First(&user, id).Error; err != nil {
			resultChan <- result{dto.UserWithProducts{}, err}
			return
		}
		resultChan <- result{dto.ConvertUserToResponse(user), nil}
	}()

	res := <-resultChan
	return res.user, res.err
}

func (u *userRepository) FindByEmail(email string) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)

	go func() {
		var user entity.User
		if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
			resultChan <- result{dto.UserWithProducts{}, err}
			return
		}
		resultChan <- result{dto.ConvertUserToResponse(user), nil}
	}()

	res := <-resultChan
	return res.user, res.err
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
