package repository

import (
	"log"
	"math"

	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(payload entity.Product) (dto.ProductWithUsers, error)
	FindByID(id uint) (dto.ProductWithUsers, error)
	FindAll(page, size int) ([]dto.ProductWithUsers, model.Paging, error)
	FindByStock(stock int) ([]dto.ProductWithUsers, error)
	UpdateByID(id uint, payload entity.Product) (dto.ProductWithUsers, error)
	DeleteByID(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// FindByStock implements ProductRepository.
func (p *productRepository) FindByStock(stock int) ([]dto.ProductWithUsers, error) {
	var products []entity.Product
	if err := p.db.Where("stock = ?", stock).Preload("Users").Find(&products).Error; err != nil {
		return nil, err
	}

	responseProducts := make([]dto.ProductWithUsers, len(products))
	for i, product := range products {
		responseProducts[i] = dto.ConvertProductToResponse(product)
	}

	return responseProducts, nil
}

// Create implements ProductRepository.
func (p *productRepository) Create(payload entity.Product) (dto.ProductWithUsers, error) {
	// Create product
	if err := p.db.Create(&payload).Error; err != nil {
		return dto.ProductWithUsers{}, err
	}

	// Load users to include in the response
	var product entity.Product
	if err := p.db.Preload("Users").First(&product, payload.ID).Error; err != nil {
		return dto.ProductWithUsers{}, err
	}

	return dto.ConvertProductToResponse(product), nil
}

// DeleteByID implements ProductRepository.
func (p *productRepository) DeleteByID(id uint) error {
	if err := p.db.Delete(&entity.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll(page int, size int) ([]dto.ProductWithUsers, model.Paging, error) {
	var products []entity.Product
	offset := (page - 1) * size

	// Retrieve total count of products
	var totalProducts int64
	if err := p.db.Model(&entity.Product{}).Count(&totalProducts).Error; err != nil {
		log.Printf("productRepository.FindAll: Error counting products: %v \n", err.Error())
		return nil, model.Paging{}, err
	}

	// Retrieve paginated products
	if err := p.db.Limit(size).Offset(offset).Preload("Users").Find(&products).Error; err != nil {
		log.Printf("productRepository.FindAll: Error fetching products: %v \n", err.Error())
		return nil, model.Paging{}, err
	}

	responseProducts := make([]dto.ProductWithUsers, len(products))
	for i, product := range products {
		responseProducts[i] = dto.ConvertProductToResponse(product)
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   int(totalProducts),
		TotalPages:  int(math.Ceil(float64(totalProducts) / float64(size))),
	}

	return responseProducts, paging, nil
}

// FindByID implements ProductRepository.
func (p *productRepository) FindByID(id uint) (dto.ProductWithUsers, error) {
	var product entity.Product
	if err := p.db.Preload("Users").First(&product, id).Error; err != nil {
		return dto.ProductWithUsers{}, err
	}
	return dto.ConvertProductToResponse(product), nil
}

// UpdateByID implements ProductRepository.
func (p *productRepository) UpdateByID(id uint, payload entity.Product) (dto.ProductWithUsers, error) {
	var product entity.Product
	if err := p.db.First(&product, id).Error; err != nil {
		return dto.ProductWithUsers{}, err
	}
	if err := p.db.Model(&product).Updates(payload).Error; err != nil {
		return dto.ProductWithUsers{}, err
	}

	p.db.Preload("Users").First(&product, id)

	return dto.ConvertProductToResponse(product), nil
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
