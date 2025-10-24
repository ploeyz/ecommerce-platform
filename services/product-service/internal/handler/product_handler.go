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

// GetProductByID handles GET /api/v1/product/:id
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
// GetAllProducts handles GET /api/v1/products
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

// UpdateProduct handles PUT /api/v1/products/:id
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

// DeleteProduct handles DELETE /api/v1/products/:id
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

// SearchProducts handles GET /api/v1/products/search
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
