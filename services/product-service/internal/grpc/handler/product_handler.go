package handler

import (
	"context"
	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"github.com/ploezy/ecommerce-platform/product-service/internal/service"
	pb "github.com/ploezy/ecommerce-platform/product-service/proto"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGRPCHandler struct {
	pb.UnimplementedProductServiceServer
	service service.ProductService
}

// NewProductGRPCHandler creates a new gRPC handler
func NewProductGRPCHandler(service service.ProductService) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		service: service,
	}
}

// CreateProduct creates a new product
func (h *ProductGRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	// Convert gRPC request to service request
	serviceReq := &model.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
		Images:      pq.StringArray(req.Images),
	}

	product, err := h.service.CreateProduct(ctx, serviceReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.ProductResponse{
		Product: h.toProtoProduct(product),
	}, nil
}

// GetProduct gets a product by ID
func (h *ProductGRPCHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.service.GetProductByID(ctx, uint(req.Id))
	if err != nil {
		if err.Error() == "product not found" {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	return &pb.ProductResponse{
		Product: h.toProtoProduct(product),
	}, nil
}

// ListProducts lists all products with pagination
func (h *ProductGRPCHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	result, err := h.service.GetAllProducts(ctx, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	// Convert to proto response
	products := make([]*pb.Product, 0)
	if productList, ok := result.Data.([]model.ProductResponse); ok {
		for _, p := range productList {
			products = append(products, h.toProtoProduct(&p))
		}
	}

	return &pb.ListProductsResponse{
		Products:   products,
		Total:      result.Total,
		Page:       int32(result.Page),
		Limit:      int32(result.Limit),
		TotalPages: int32(result.TotalPages),
	}, nil
}

// UpdateProduct updates a product
func (h *ProductGRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	serviceReq := &model.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		Category:    req.Category,
		Images:      pq.StringArray(req.Images),
	}

	product, err := h.service.UpdateProduct(ctx, uint(req.Id), serviceReq)
	if err != nil {
		if err.Error() == "product not found" {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &pb.ProductResponse{
		Product: h.toProtoProduct(product),
	}, nil
}

// DeleteProduct deletes a product
func (h *ProductGRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.service.DeleteProduct(ctx, uint(req.Id))
	if err != nil {
		if err.Error() == "product not found" {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &pb.DeleteProductResponse{
		Success: true,
		Message: "Product deleted successfully",
	}, nil
}

// SearchProducts searches products by keyword
func (h *ProductGRPCHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	result, err := h.service.SearchProducts(ctx, req.Keyword, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search products: %v", err)
	}

	// Convert to proto response
	products := make([]*pb.Product, 0)
	if productList, ok := result.Data.([]model.ProductResponse); ok {
		for _, p := range productList {
			products = append(products, h.toProtoProduct(&p))
		}
	}

	return &pb.ListProductsResponse{
		Products:   products,
		Total:      result.Total,
		Page:       int32(result.Page),
		Limit:      int32(result.Limit),
		TotalPages: int32(result.TotalPages),
	}, nil
}

// Helper function to convert ProductResponse to Proto Product
func (h *ProductGRPCHandler) toProtoProduct(p *model.ProductResponse) *pb.Product {
	return &pb.Product{
		Id:          uint32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		Category:    p.Category,
		Images:      p.Images,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}