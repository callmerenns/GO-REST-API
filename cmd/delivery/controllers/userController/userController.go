package userController

import (
	"log"
	"net/http"
	"strconv"

	"github.com/altsaqif/go-rest/cmd/config"
	"github.com/altsaqif/go-rest/cmd/delivery/middlewares"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/common"
	"github.com/altsaqif/go-rest/cmd/shared/model"
	"github.com/altsaqif/go-rest/cmd/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUc  usecase.UserUseCase
	rg      *gin.RouterGroup
	authMid middlewares.AuthMiddleware
}

// @Summary Get all users
// @Description Get a list of all users with pagination
// @Tags users
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Status
// @Router /profiles [get]
func (u *UserController) GetAllHandler(ctx *gin.Context) {
	page, errPage := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if errPage != nil {
		log.Fatal("Error : ", errPage)
	}

	size, errSize := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if errSize != nil {
		log.Fatal("Error : ", errSize)
	}

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	type result struct {
		users  []dto.UserWithProducts
		paging model.Paging
		err    error
	}

	resultChan := make(chan result)
	go func() {
		users, paging, err := u.userUc.FindAllUsers(page, size)
		resultChan <- result{users, paging, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	var interfaceSlice = make([]interface{}, len(res.users))
	for i, v := range res.users {
		interfaceSlice[i] = v
	}

	common.SendPagedResponse(ctx, interfaceSlice, res.paging, "Ok")
}

// @Summary Get user by ID
// @Description Get details of a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Router /profiles/{id} [get]
func (u *UserController) GetHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	convUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing string to uint64 : %v", err)
	}

	uintValue := uint(convUint)

	type result struct {
		user dto.UserWithProducts
		err  error
	}

	resultChan := make(chan result)
	go func() {
		user, err := u.userUc.FindUserByID(uintValue)
		resultChan <- result{user, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, res.err.Error())
		return
	}
	common.SendSingleResponse(ctx, "Ok", res.user)
}

func (u *UserController) Route() {
	u.rg.GET(config.GetUsersList, u.authMid.RequireToken("user", "admin"), u.GetAllHandler)
	u.rg.GET(config.GetUsers, u.authMid.RequireToken("user", "admin"), u.GetHandler)
}

func NewUserController(userUc usecase.UserUseCase, rg *gin.RouterGroup, authMid middlewares.AuthMiddleware) *UserController {
	return &UserController{userUc: userUc, rg: rg, authMid: authMid}
}
