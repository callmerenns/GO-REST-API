package usecase

import (
	"fmt"
	"time"

	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/repository"
	"github.com/altsaqif/go-rest/cmd/shared/model"
)

type UserUseCase interface {
	RegisterNewUser(payload entity.User) (dto.UserWithProducts, error)
	FindUserByID(id uint) (dto.UserWithProducts, error)
	FindUserByEmail(email string) (dto.UserWithProducts, error)
	FindAllUsers(page, size int) ([]dto.UserWithProducts, model.Paging, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

// FindAllUsers implements UserUseCase.
func (u *userUseCase) FindAllUsers(page, size int) ([]dto.UserWithProducts, model.Paging, error) {
	return u.repo.FindAll(page, size)
}

// GetUserByEmail implements UserUseCase.
func (u *userUseCase) FindUserByEmail(email string) (dto.UserWithProducts, error) {
	return u.repo.FindByEmail(email)
}

func (u *userUseCase) RegisterNewUser(payload entity.User) (dto.UserWithProducts, error) {
	userExist, _ := u.repo.FindByEmail(payload.Email)
	if userExist.Email == payload.Email {
		return dto.UserWithProducts{}, fmt.Errorf("user with email: %s already exists", payload.Email)
	}
	payload.UpdatedAt = time.Now()
	createdUser, err := u.repo.Create(payload)
	if err != nil {
		return dto.UserWithProducts{}, err
	}
	return createdUser, nil
}

func (u *userUseCase) FindUserByID(id uint) (dto.UserWithProducts, error) {
	return u.repo.FindByID(id)
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
