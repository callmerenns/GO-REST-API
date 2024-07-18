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
	type result struct {
		products []entity.Product
		err      error
	}

	resultChan := make(chan result)
	go func() {
		var products []entity.Product
		err := p.db.Where("stock = ?", stock).Preload("Users").Find(&products).Error
		resultChan <- result{products, err}
	}()

	res := <-resultChan
	if res.err != nil {
		return nil, res.err
	}

	responseProducts := make([]dto.ProductWithUsers, len(res.products))
	for i, product := range res.products {
		responseProducts[i] = dto.ConvertProductToResponse(product)
	}

	return responseProducts, nil
}

// Create implements ProductRepository.
func (p *productRepository) Create(payload entity.Product) (dto.ProductWithUsers, error) {
	type result struct {
		product entity.Product
		err     error
	}

	resultChan := make(chan result)
	go func() {
		if err := p.db.Create(&payload).Error; err != nil {
			resultChan <- result{entity.Product{}, err}
			return
		}

		var product entity.Product
		err := p.db.Preload("Users").First(&product, payload.ID).Error
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	if res.err != nil {
		return dto.ProductWithUsers{}, res.err
	}

	return dto.ConvertProductToResponse(res.product), nil
}

// DeleteByID implements ProductRepository.
func (p *productRepository) DeleteByID(id uint) error {
	type result struct {
		err error
	}

	resultChan := make(chan result)
	go func() {
		err := p.db.Delete(&entity.Product{}, id).Error
		resultChan <- result{err}
	}()

	res := <-resultChan
	return res.err
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll(page int, size int) ([]dto.ProductWithUsers, model.Paging, error) {
	type result struct {
		totalProducts int64
		products      []entity.Product
		err           error
	}

	offset := (page - 1) * size
	resultChan := make(chan result)

	go func() {
		var totalProducts int64
		if err := p.db.Model(&entity.Product{}).Count(&totalProducts).Error; err != nil {
			resultChan <- result{0, nil, err}
			return
		}

		var products []entity.Product
		if err := p.db.Limit(size).Offset(offset).Preload("Users").Find(&products).Error; err != nil {
			resultChan <- result{totalProducts, nil, err}
			return
		}

		resultChan <- result{totalProducts, products, nil}
	}()

	res := <-resultChan
	if res.err != nil {
		log.Printf("productRepository.FindAll: Error: %v \n", res.err)
		return nil, model.Paging{}, res.err
	}

	responseProducts := make([]dto.ProductWithUsers, len(res.products))
	for i, product := range res.products {
		responseProducts[i] = dto.ConvertProductToResponse(product)
	}

	paging := model.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   int(res.totalProducts),
		TotalPages:  int(math.Ceil(float64(res.totalProducts) / float64(size))),
	}

	return responseProducts, paging, nil
}

// FindByID implements ProductRepository.
func (p *productRepository) FindByID(id uint) (dto.ProductWithUsers, error) {
	type result struct {
		product entity.Product
		err     error
	}

	resultChan := make(chan result)
	go func() {
		var product entity.Product
		err := p.db.Preload("Users").First(&product, id).Error
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	if res.err != nil {
		return dto.ProductWithUsers{}, res.err
	}

	return dto.ConvertProductToResponse(res.product), nil
}

// UpdateByID implements ProductRepository.
func (p *productRepository) UpdateByID(id uint, payload entity.Product) (dto.ProductWithUsers, error) {
	type result struct {
		product entity.Product
		err     error
	}

	resultChan := make(chan result)
	go func() {
		var product entity.Product
		if err := p.db.First(&product, id).Error; err != nil {
			resultChan <- result{entity.Product{}, err}
			return
		}
		if err := p.db.Model(&product).Updates(payload).Error; err != nil {
			resultChan <- result{entity.Product{}, err}
			return
		}

		p.db.Preload("Users").First(&product, id)
		resultChan <- result{product, nil}
	}()

	res := <-resultChan
	if res.err != nil {
		return dto.ProductWithUsers{}, res.err
	}

	return dto.ConvertProductToResponse(res.product), nil
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
