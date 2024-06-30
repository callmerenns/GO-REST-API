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
	return a.uc.FindUserByEmail(email)
}

func (a *authUseCase) Login(payload dto.AuthRequestLoginDto) (dto.AuthResponseDto, error) {
	user, err := a.uc.FindUserByEmail(payload.Email)
	if err != nil || !utils.CheckPasswordHash(payload.Password, user.Password) {
		return dto.AuthResponseDto{}, err
	}

	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}
	return token, nil
}

func (a *authUseCase) Register(payload dto.AuthRequestRegisterDto) (dto.UserWithProducts, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return dto.UserWithProducts{}, err
	}

	return a.uc.RegisterNewUser(entity.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
		Role:      payload.Role,
	})
}

func NewAuthUseCase(uc UserUseCase, jwtService service.JwtService) AuthUseCase {
	return &authUseCase{uc: uc, jwtService: jwtService}
}
