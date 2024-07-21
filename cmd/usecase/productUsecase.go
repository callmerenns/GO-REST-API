package usecase

import (
	"fmt"

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
	ProductExists(id uint) (bool, error)
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (p *productUseCase) CreateProduct(payload entity.Product) (dto.ProductWithUsers, error) {
	type result struct {
		product dto.ProductWithUsers
		err     error
	}

	resultChan := make(chan result)
	go func() {
		product, err := p.repo.Create(payload)
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	return res.product, res.err
}

func (p *productUseCase) FindProductByID(id uint) (dto.ProductWithUsers, error) {
	type result struct {
		product dto.ProductWithUsers
		err     error
	}

	resultChan := make(chan result)
	go func() {
		product, err := p.repo.FindByID(id)
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	return res.product, res.err
}

func (p *productUseCase) FindAllProducts(page, size int) ([]dto.ProductWithUsers, model.Paging, error) {
	type result struct {
		products []dto.ProductWithUsers
		paging   model.Paging
		err      error
	}

	resultChan := make(chan result)
	go func() {
		products, paging, err := p.repo.FindAll(page, size)
		resultChan <- result{products, paging, err}
	}()

	res := <-resultChan
	return res.products, res.paging, res.err
}

// FindProductsByStock implements ProductUseCase.
func (p *productUseCase) FindProductsByStock(stock int) ([]dto.ProductWithUsers, error) {
	type result struct {
		products []dto.ProductWithUsers
		err      error
	}

	resultChan := make(chan result)
	go func() {
		products, err := p.repo.FindByStock(stock)
		resultChan <- result{products, err}
	}()

	res := <-resultChan
	return res.products, res.err
}

func (p *productUseCase) UpdateProduct(id uint, payload entity.Product) (dto.ProductWithUsers, error) {
	type result struct {
		product dto.ProductWithUsers
		err     error
	}

	resultChan := make(chan result)
	go func() {
		product, err := p.repo.UpdateByID(id, payload)
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	return res.product, res.err
}

func (p *productUseCase) DeleteProduct(id uint) error {
	type result struct {
		err error
	}

	resultChan := make(chan result)
	go func() {
		err := p.repo.DeleteByID(id)
		resultChan <- result{err}
	}()

	res := <-resultChan
	return res.err
}

func (p *productUseCase) ProductExists(id uint) (bool, error) {
	// Channels for signaling completion and errors
	existsCh := make(chan bool)
	errCh := make(chan error, 1)

	// Goroutine to check if the product exists
	go func() {
		exists, err := p.repo.ProductExists(id)
		if err != nil {
			errCh <- fmt.Errorf("failed to check if product exists: %v", err)
			return
		}
		existsCh <- exists
	}()

	// Wait for results or errors
	select {
	case exists := <-existsCh:
		return exists, nil
	case err := <-errCh:
		return false, err
	}
}
