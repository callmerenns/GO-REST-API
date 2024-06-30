package usecase

import (
	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/repository"
	"github.com/altsaqif/go-rest/cmd/shared/model"
)

type ProductUseCase interface {
	CreateProduct(payload entity.Product) (dto.ProductWithUsers, error)
	FindProductByID(id uint) (dto.ProductWithUsers, error)
	FindAllProducts(page, size int) ([]dto.ProductWithUsers, model.Paging, error)
	FindProductsByStock(stock int) ([]dto.ProductWithUsers, error)
	UpdateProduct(id uint, payload entity.Product) (dto.ProductWithUsers, error)
	DeleteProduct(id uint) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

// FindProductsByStock implements ProductUseCase.
func (p *productUseCase) FindProductsByStock(stock int) ([]dto.ProductWithUsers, error) {
	return p.repo.FindByStock(stock)
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (p *productUseCase) CreateProduct(payload entity.Product) (dto.ProductWithUsers, error) {
	return p.repo.Create(payload)
}

func (p *productUseCase) FindProductByID(id uint) (dto.ProductWithUsers, error) {
	return p.repo.FindByID(id)
}

func (p *productUseCase) FindAllProducts(page, size int) ([]dto.ProductWithUsers, model.Paging, error) {
	return p.repo.FindAll(page, size)
}

func (p *productUseCase) UpdateProduct(id uint, payload entity.Product) (dto.ProductWithUsers, error) {
	return p.repo.UpdateByID(id, payload)
}

func (p *productUseCase) DeleteProduct(id uint) error {
	return p.repo.DeleteByID(id)
}
