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
	type result struct {
		users  []dto.UserWithProducts
		paging model.Paging
		err    error
	}

	resultChan := make(chan result)
	go func() {
		users, paging, err := u.repo.FindAll(page, size)
		resultChan <- result{users, paging, err}
	}()

	res := <-resultChan
	return res.users, res.paging, res.err
}

// GetUserByEmail implements UserUseCase.
func (u *userUseCase) FindUserByEmail(email string) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := u.repo.FindByEmail(email)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	return res.user, res.err
}

func (u *userUseCase) RegisterNewUser(payload entity.User) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		userExist, _ := u.repo.FindByEmail(payload.Email)
		if userExist.Email == payload.Email {
			resultChan <- result{dto.UserWithProducts{}, fmt.Errorf("user with email: %s already exists", payload.Email)}
			return
		}
		payload.UpdatedAt = time.Now()
		createdUser, err := u.repo.Create(payload)
		if err != nil {
			resultChan <- result{dto.UserWithProducts{}, err}
			return
		}
		resultChan <- result{createdUser, nil}
	}()

	res := <-resultChan
	return res.user, res.err
}

func (u *userUseCase) FindUserByID(id uint) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := u.repo.FindByID(id)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	return res.user, res.err
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
