package productController

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/altsaqif/go-rest/cmd/config"
	"github.com/altsaqif/go-rest/cmd/delivery/middlewares"
	"github.com/altsaqif/go-rest/cmd/entity"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/common"
	"github.com/altsaqif/go-rest/cmd/shared/model"
	"github.com/altsaqif/go-rest/cmd/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	productUc usecase.ProductUseCase
	rg        *gin.RouterGroup
	authMid   middlewares.AuthMiddleware
}

func NewProductController(productUc usecase.ProductUseCase, rg *gin.RouterGroup, authMid middlewares.AuthMiddleware) *ProductController {
	return &ProductController{productUc: productUc, rg: rg, authMid: authMid}
}

// @Summary Get all products
// @Description Get a list of all products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Status
// @Failure 404 {object} model.Status
// @Router /products [get]
func (p *ProductController) GetAllHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	type result struct {
		products []dto.ProductWithUsers
		paging   model.Paging
		err      error
	}

	resultChan := make(chan result)
	go func() {
		products, paging, err := p.productUc.FindAllProducts(page, size)
		resultChan <- result{products, paging, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	if len(res.products) == 0 {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Products not found")
		return
	}

	var interfaceSlice = make([]interface{}, len(res.products))
	for i, v := range res.products {
		interfaceSlice[i] = v
	}

	common.SendPagedResponse(ctx, interfaceSlice, res.paging, "Ok")
}

// @Summary Get product by ID
// @Description Get details of a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 404 {object} model.Status
// @Failure 500 {object} model.Status
// @Failure 404 {object} model.Status
// @Router /products/{id} [get]
func (p *ProductController) GetByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	convUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	uintValue := uint(convUint)
	type result struct {
		product dto.ProductWithUsers
		err     error
	}

	resultChan := make(chan result)
	go func() {
		product, err := p.productUc.FindProductByID(uintValue)
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	if res.err != nil {
		// Customize the error message for "record not found"
		if errors.Is(res.err, gorm.ErrRecordNotFound) {
			common.SendErrorResponse(ctx, http.StatusNotFound, "Product not found")
		} else {
			common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		}
		return
	}

	// Check if the product is empty
	if res.product.ID == 0 {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Product not found")
		return
	}

	common.SendSingleResponse(ctx, "Ok", res.product)
}

// @Summary Get products by stock
// @Description Get a list of products by stock value
// @Tags products
// @Produce json
// @Param stock path int true "Stock"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 500 {object} model.Status
// @Failure 404 {object} model.Status
// @Router /products/stock/{stock} [get]
func (p *ProductController) GetByStockHandler(ctx *gin.Context) {
	stockStr := ctx.Param("stock")
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid stock value")
		return
	}

	type result struct {
		products []dto.ProductWithUsers
		err      error
	}

	resultChan := make(chan result)
	go func() {
		products, err := p.productUc.FindProductsByStock(stock)
		resultChan <- result{products, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	// Check if no products are found
	if len(res.products) == 0 {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Products not found")
		return
	}

	common.SendSingleResponse(ctx, "Ok", res.products)
}

// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param Product body entity.Product true "Product Payload"
// @Success 201 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 500 {object} model.Status
// @Router /products [post]
func (p *ProductController) CreateHandler(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	type result struct {
		createdProduct dto.ProductWithUsers
		err            error
	}

	resultChan := make(chan result)
	go func() {
		createdProduct, err := p.productUc.CreateProduct(product)
		resultChan <- result{createdProduct, err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	common.SendSingleResponse(ctx, "Product created successfully", res.createdProduct)
}

// @Summary Update product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param Product body entity.Product true "Product Payload"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 400 {object} model.Status
// @Failure 404 {object} model.Status
// @Failure 500 {object} model.Status
// @Router /products/{id} [put]
func (p *ProductController) UpdateHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	convUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	uintValue := uint(convUint)
	var payload entity.Product
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	type result struct {
		product dto.ProductWithUsers
		err     error
	}

	resultChan := make(chan result)
	go func() {
		product, err := p.productUc.UpdateProduct(uintValue, payload)
		resultChan <- result{product, err}
	}()

	res := <-resultChan
	if res.err != nil {
		if errors.Is(res.err, gorm.ErrRecordNotFound) {
			common.SendErrorResponse(ctx, http.StatusNotFound, "Product not found")
		} else {
			common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		}
		return
	}
	common.SendSingleResponse(ctx, "Product updated successfully", res.product)
}

// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} model.SingleResponse
// @Failure 400 {object} model.Status
// @Failure 500 {object} model.Status
// @Failure 404 {object} model.Status
// @Failure 500 {object} model.Status
// @Router /products/{id} [delete]
func (p *ProductController) DeleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	convUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	uintValue := uint(convUint)

	// First, check if the product exists
	exists, err := p.productUc.ProductExists(uintValue)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if !exists {
		common.SendErrorResponse(ctx, http.StatusNotFound, "Product not found")
		return
	}

	// Proceed to delete the product
	type result struct {
		err error
	}

	resultChan := make(chan result)
	go func() {
		err := p.productUc.DeleteProduct(uintValue)
		resultChan <- result{err}
	}()

	res := <-resultChan
	if res.err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, res.err.Error())
		return
	}

	common.SendSuccessResponse(ctx, "Product deleted successfully")
}

func (p *ProductController) Route() {
	p.rg.GET(config.GetProductsList, p.authMid.RequireToken("customer", "reseller", "admin"), p.GetAllHandler)
	p.rg.GET(config.GetProducts, p.authMid.RequireToken("customer", "reseller", "admin"), p.GetByIDHandler)
	p.rg.GET(config.GetProductsByStocks, p.authMid.RequireToken("reseller", "admin"), p.GetByStockHandler)
	p.rg.POST(config.PostProducts, p.authMid.RequireToken("reseller", "admin"), p.CreateHandler)
	p.rg.PUT(config.PutProducts, p.authMid.RequireToken("reseller", "admin"), p.UpdateHandler)
	p.rg.DELETE(config.DelProducts, p.authMid.RequireToken("reseller", "admin"), p.DeleteHandler)
}
