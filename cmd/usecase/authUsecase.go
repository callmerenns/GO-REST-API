package usecase

import (
	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/service"
	"github.com/altsaqif/go-rest/cmd/utils"
)

type AuthUseCase interface {
	Login(payload dto.AuthRequestLoginDto) (dto.AuthResponseDto, error)
	Register(payload dto.AuthRequestRegisterDto) (dto.UserWithProducts, error)
	FindUserByEmail(email string) (dto.UserWithProducts, error)
}

type authUseCase struct {
	uc         UserUseCase
	jwtService service.JwtService
}

// GetUserByEmail implements AuthUseCase.
func (a *authUseCase) FindUserByEmail(email string) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := a.uc.FindUserByEmail(email)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	return res.user, res.err
}

func (a *authUseCase) Login(payload dto.AuthRequestLoginDto) (dto.AuthResponseDto, error) {
	type result struct {
		token dto.AuthResponseDto
		err   error
	}

	resultChan := make(chan result)
	go func() {
		user, err := a.uc.FindUserByEmail(payload.Email)
		if err != nil || !utils.CheckPasswordHash(payload.Password, user.Password) {
			resultChan <- result{dto.AuthResponseDto{}, err}
			return
		}

		token, err := a.jwtService.CreateToken(user)
		resultChan <- result{token, err}
	}()

	res := <-resultChan
	return res.token, res.err
}

func (a *authUseCase) Register(payload dto.AuthRequestRegisterDto) (dto.UserWithProducts, error) {
	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			resultChan <- result{dto.UserWithProducts{}, err}
			return
		}

		user, err := a.uc.RegisterNewUser(entity.User{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  hashedPassword,
			Role:      payload.Role,
		})
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	return res.user, res.err
}

func NewAuthUseCase(uc UserUseCase, jwtService service.JwtService) AuthUseCase {
	return &authUseCase{uc: uc, jwtService: jwtService}
}
