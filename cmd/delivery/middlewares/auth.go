package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/altsaqif/go-rest/cmd/shared/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authHeader AuthHeader
		if err := ctx.ShouldBindHeader(&authHeader); err != nil {
			log.Printf("RequireToken: Error binding header: %v \n", err)
		}

		tokenHeader := strings.TrimPrefix(authHeader.AuthorizationHeader, "Bearer ")

		if tokenHeader == "" {
			// Log when checking the cookie
			log.Println("RequireToken: Checking cookie for token")
			cookie, err := ctx.Cookie("token")
			fmt.Println("Token : ", cookie)
			if err != nil {
				log.Println("RequireToken: Error retrieving token from cookie:", err)
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			tokenHeader = cookie
			log.Printf("RequireToken: Retrieved token from cookie: %v\n", tokenHeader)
		}

		if tokenHeader == "" {
			log.Println("RequireToken: Token is empty")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := a.jwtService.ParseToken(tokenHeader)
		if err != nil {
			log.Printf("RequireToken: Error parsing token: %v \n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", claims["userId"])

		role, ok := claims["role"]
		if !ok {
			log.Println("RequireToken: Missing role in token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !isValidRole(role.(string), roles) {
			log.Println("RequireToken: Invalid role")
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

func isValidRole(userRole string, validRoles []string) bool {
	for _, role := range validRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
