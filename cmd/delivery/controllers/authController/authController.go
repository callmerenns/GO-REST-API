package authController

import (
	"net/http"

	"github.com/altsaqif/go-rest/cmd/config"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/common"
	"github.com/altsaqif/go-rest/cmd/usecase"
	"github.com/altsaqif/go-rest/cmd/utils"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication
type AuthController struct {
	authUc usecase.AuthUseCase
	rg     *gin.RouterGroup
}

// NewAuthController creates a new AuthController
func NewAuthController(authUc usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUc: authUc, rg: rg}
}

// @Summary Login user
// @Description Log in an existing user
// @Tags auth
// @Accept json
// @Produce json
// @Param AuthRequestLoginDto body dto.AuthRequestLoginDto true "Login Payload"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 401 {object} model.Status
// @Failure 500 {object} model.Status
// @Router /auth/login [post]
func (a *AuthController) loginHandler(ctx *gin.Context) {
	var payload dto.AuthRequestLoginDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := a.authUc.FindUserByEmail(payload.Email)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	if res.err != nil || !utils.CheckPasswordHash(payload.Password, res.user.Password) {
		common.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := a.authUc.Login(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Set token to cookie
	ctx.SetCookie("token", token.Token, 3600, "/", "", false, true)

	common.SendSuccessResponse(ctx, "Successfully Login")
}

// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param AuthRequestRegisterDto body dto.AuthRequestRegisterDto true "Register Payload"
// @Success 201 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 400 {object} model.Status
// @Failure 500 {object} model.Status
// @Router /auth/register [post]
func (a *AuthController) registerHandler(ctx *gin.Context) {
	var payload dto.AuthRequestRegisterDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Password != payload.PasswordConfirm {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Password not match")
		return
	}

	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := a.authUc.Register(payload)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	responseData := map[string]interface{}{
		"id":         res.user.ID,
		"username":   res.user.FirstName + " " + res.user.LastName,
		"email":      res.user.Email,
		"role":       res.user.Role,
		"created_at": res.user.CreatedAt,
		"updated_at": res.user.UpdatedAt,
		"deleted_at": res.user.DeletedAt,
	}

	common.SendCreateResponse(ctx, "User registered successfully", responseData)
}

// @Summary Logout user
// @Description Log out the current user
// @Tags auth
// @Produce json
// @Success 200 {object} model.SingleResponse
// @Router /auth/logout [get]
func (a *AuthController) logoutHandler(ctx *gin.Context) {
	// Clear the token cookie
	ctx.SetCookie("token", "", -1, "/", "", false, true)

	common.SendSuccessResponse(ctx, "Logout successfully!")
}

// Route initializes the auth routes
func (a *AuthController) Route() {
	a.rg.POST(config.PostLogin, a.loginHandler)
	a.rg.POST(config.PostRegister, a.registerHandler)
	a.rg.GET(config.GetLogout, a.logoutHandler)
}
