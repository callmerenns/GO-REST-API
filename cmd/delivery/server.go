package delivery

import (
	"fmt"
	"log"

	"github.com/altsaqif/go-rest/cmd/config"
	"github.com/altsaqif/go-rest/cmd/delivery/controllers/authController"
	"github.com/altsaqif/go-rest/cmd/delivery/controllers/productController"
	"github.com/altsaqif/go-rest/cmd/delivery/controllers/userController"
	"github.com/altsaqif/go-rest/cmd/delivery/middlewares"
	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/repository"
	"github.com/altsaqif/go-rest/cmd/shared/service"
	"github.com/altsaqif/go-rest/cmd/usecase"
	_ "github.com/altsaqif/go-rest/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	productUc  usecase.ProductUseCase
	userUc     usecase.UserUseCase
	authUc     usecase.AuthUseCase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.ApiGroup)
	authMid := middlewares.NewAuthMiddleware(s.jwtService)
	authController.NewAuthController(s.authUc, rg).Route()
	userController.NewUserController(s.userUc, rg, authMid).Route()
	productController.NewProductController(s.productUc, rg, authMid).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connection error")
	}

	// Drop existing tables
	db.Migrator().DropTable(&entity.User{}, &entity.Product{}, &entity.Enrollment{})

	err = db.AutoMigrate(&entity.User{}, &entity.Product{}, &entity.Enrollment{})

	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db.SetupJoinTable(&entity.User{}, "Products", &entity.Enrollment{})

	jwtService := service.NewJwtService(cfg.TokenConfig)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	productUc := usecase.NewProductUseCase(productRepo)
	userUc := usecase.NewUserUseCase(userRepo)
	authUc := usecase.NewAuthUseCase(userUc, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	// Swagger handler
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		productUc:  productUc,
		userUc:     userUc,
		authUc:     authUc,
		jwtService: jwtService,
		engine:     engine,
		host:       host,
	}
}
