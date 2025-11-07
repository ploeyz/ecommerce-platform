package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"github.com/ploezy/ecommerce-platform/product-service/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}
// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequest true "Product Data"
// @Success 201 {object} Response{data=model.ProductResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Security BearerAuth
// @Router /products [post]
// CreateProduct handler POST /api/v1/products
func (h * ProductHandler) CreateProduct(c *gin.Context){
	var req model.CreateProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	product, err := h.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		ErrorResponse(c,http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated,"Product created successfully",product)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a single product by ID (with Redis cache)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Response{data=model.ProductResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /products/{id} [get]
func (h * ProductHandler) GetProductByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr,10,32)
	if err != nil {
		ErrorResponse(c,http.StatusBadRequest, "Invaild product ID")
		return
	}

	product, err := h.service.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK,"Product retrieved successfully",product)
}
// GetAllProducts godoc
// @Summary Get all products
// @Description Get all products with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} Response{data=model.PaginationResponse}
// @Failure 500 {object} Response
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.service.GetAllProducts(c.Request.Context(), page, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.UpdateProductRequest true "Product Data"
// @Success 200 {object} Response{data=model.ProductResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Security BearerAuth
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "product not found" {
			ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product (Admin only, soft delete)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "product not found" {
			ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products by keyword in name or category
// @Tags Products
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} Response{data=model.PaginationResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if keyword == "" {
		ErrorResponse(c, http.StatusBadRequest, "keyword is required")
		return
	}

	products, err := h.service.SearchProducts(c.Request.Context(), keyword, page, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Search completed successfully", products)
}
