// shared/common/jsonResponse.go

package common

import (
	"net/http"

	"github.com/altsaqif/go-rest/cmd/shared/model"
	"github.com/gin-gonic/gin"
)

// SendCreateResponse defines the standard create response structure
func SendCreateResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, &model.SingleResponse{
		Status: model.Status{
			Code:    http.StatusCreated,
			Message: message,
		},
		Data: data,
	})
}

// SendSuccessResponse defines the standard success response structure
func SendSuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &model.SingleResponse{
		Status: model.Status{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Data: data,
	})
}

// SendSingleResponse defines the standard single response structure
func SendSingleResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, &model.SingleResponse{
		Status: model.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

// SendPagedResponse defines the standard paged response structure
func SendPagedResponse(ctx *gin.Context, data []interface{}, paging model.Paging, message string) {
	ctx.JSON(http.StatusOK, &model.PagedResponse{
		Status: model.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

// SendErrorResponse defines the standard error response structure
func SendErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, &model.Status{
		Code:    code,
		Message: message,
	})
}
