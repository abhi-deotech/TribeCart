package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tribecart/proto/tribecart/v1"
	"github.com/tribecart/services/products/internal/repository"
)

var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = status.Error(codes.NotFound, "product not found")
	// ErrInvalidArgument is returned when an invalid argument is provided
	ErrInvalidArgument = status.Error(codes.InvalidArgument, "invalid argument")
)

// ProductService implements the ProductServiceServer interface
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Product, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be greater than 0")
	}
	if req.StockQuantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock quantity cannot be negative")
	}

	// Map request to product
	product := &v1.Product{
		Id:              uuid.New().String(),
		SellerId:        req.SellerId,
		Name:            req.Name,
		Description:     req.Description,
		Sku:             req.Sku,
		Barcode:         req.Barcode,
		Price:           req.Price,
		SalePrice:       req.SalePrice,
		CostPrice:       req.CostPrice,
		StockQuantity:   req.StockQuantity,
		TrackInventory:  req.TrackInventory,
		MinStockLevel:   req.MinStockLevel,
		Weight:          req.Weight,
		Length:          req.Length,
		Width:           req.Width,
		Height:          req.Height,
		Type:            req.Type,
		Status:          v1.ProductStatus_PRODUCT_STATUS_DRAFT,
		IsFeatured:      req.IsFeatured,
		IsVisible:       req.IsVisible,
		RequiresShipping: req.RequiresShipping,
		IsTaxable:       req.IsTaxable,
		TaxClassId:      req.TaxClassId,
		CategoryIds:     req.CategoryIds,
		Tags:            req.Tags,
		Images:          req.Images,
		Specifications:  req.Specifications,
		SeoTitle:        req.SeoTitle,
		SeoDescription:  req.SeoDescription,
		SeoKeywords:     req.SeoKeywords,
	}

	// Set default status if not provided
	if req.Status != v1.ProductStatus_PRODUCT_STATUS_UNSPECIFIED {
		product.Status = req.Status
	} else {
		product.Status = v1.ProductStatus_PRODUCT_STATUS_ACTIVE
	}

	// Save to repository
	createdProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	// TODO: Handle variants if provided

	return createdProduct, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
	if req.Id == "" {
		return nil, ErrInvalidArgument
	}

	product, err := s.repo.GetProduct(ctx, req.Id, req.IncludeDeleted)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	// TODO: Load variants if include_variants is true

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Product, error) {
	if req.Id == "" {
		return nil, ErrInvalidArgument
	}

	// Get existing product
	existing, err := s.repo.GetProduct(ctx, req.Id, false)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	if existing == nil {
		return nil, ErrProductNotFound
	}

	// Apply field mask if provided
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// In a real implementation, you would use the field mask to update only the specified fields
		// For simplicity, we'll just update all fields here
	}

	// Update fields that can be updated
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.StockQuantity >= 0 {
		existing.StockQuantity = req.StockQuantity
	}
	if req.Sku != "" {
		existing.Sku = req.Sku
	}
	if req.Barcode != "" {
		existing.Barcode = req.Barcode
	}
	if req.SalePrice >= 0 {
		existing.SalePrice = req.SalePrice
	}
	if req.CostPrice >= 0 {
		existing.CostPrice = req.CostPrice
	}
	if req.MinStockLevel >= 0 {
		existing.MinStockLevel = req.MinStockLevel
	}
	if req.Weight >= 0 {
		existing.Weight = req.Weight
	}
	if req.Length >= 0 {
		existing.Length = req.Length
	}
	if req.Width >= 0 {
		existing.Width = req.Width
	}
	if req.Height >= 0 {
		existing.Height = req.Height
	}
	if req.Type != v1.ProductType_PRODUCT_TYPE_UNSPECIFIED {
		existing.Type = req.Type
	}
	if req.Status != v1.ProductStatus_PRODUCT_STATUS_UNSPECIFIED {
		existing.Status = req.Status
	}

	existing.IsFeatured = req.IsFeatured
	existing.IsVisible = req.IsVisible
	existing.RequiresShipping = req.RequiresShipping
	existing.IsTaxable = req.IsTaxable

	if req.TaxClassId != "" {
		existing.TaxClassId = req.TaxClassId
	}

	if req.CategoryIds != nil {
		existing.CategoryIds = req.CategoryIds
	}

	if req.Tags != nil {
		existing.Tags = req.Tags
	}

	if req.Images != nil {
		existing.Images = req.Images
	}

	if req.Specifications != nil {
		existing.Specifications = req.Specifications
	}

	if req.SeoTitle != "" {
		existing.SeoTitle = req.SeoTitle
	}

	if req.SeoDescription != "" {
		existing.SeoDescription = req.SeoDescription
	}

	if req.SeoKeywords != "" {
		existing.SeoKeywords = req.SeoKeywords
	}

	// Update in repository
	updatedProduct, err := s.repo.UpdateProduct(ctx, existing, nil) // Pass nil for updateMask for now
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	// TODO: Handle variants if provided

	return updatedProduct, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, ErrInvalidArgument
	}

	err := s.repo.DeleteProduct(ctx, req.Id, req.Force)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListProducts retrieves a list of products with filtering and pagination
func (s *ProductService) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	// Set default pagination values
	if req.PageSize <= 0 {
		req.PageSize = 20 // Default page size
	}
	if req.Page <= 0 {
		req.Page = 1 // Default to first page
	}

	// Get products from repository
	products, totalCount, err := s.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	return &v1.ListProductsResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// UpdateStock updates the stock quantity for a product or variant
func (s *ProductService) UpdateStock(ctx context.Context, req *v1.UpdateStockRequest) (*v1.Product, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
	if req.Quantity == 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity cannot be zero")
	}

	// Update stock in repository
	product, err := s.repo.UpdateStock(ctx, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}

	return product, nil
}

// SearchProducts searches for products based on a query
func (s *ProductService) SearchProducts(ctx context.Context, req *v1.SearchProductsRequest) (*v1.SearchProductsResponse, error) {
	// For now, we'll just forward this to ListProducts with the search query
	// In a real implementation, you would use a full-text search engine like Elasticsearch

	listReq := &v1.ListProductsRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		SearchQuery:   req.Query,
		CategoryIds:   req.CategoryIds,
		MinPrice:      req.MinPrice,
		MaxPrice:      req.MaxPrice,
		IncludeDeleted: req.IncludeOutOfStock,
	}

	products, totalCount, err := s.repo.ListProducts(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search products: %v", err)
	}

	return &v1.SearchProductsResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// GetProductsByCategory retrieves products by category
func (s *ProductService) GetProductsByCategory(ctx context.Context, req *v1.GetProductsByCategoryRequest) (*v1.ListProductsResponse, error) {
	if req.CategoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "category_id is required")
	}

	// Set default pagination values
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// Create list request with category filter
	listReq := &v1.ListProductsRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		CategoryIds:   []string{req.CategoryId},
		IncludeDeleted: false,
	}

	// Get products from repository
	products, totalCount, err := s.repo.ListProducts(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products by category: %v", err)
	}

	return &v1.ListProductsResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// BulkImportProducts imports multiple products from a file
func (s *ProductService) BulkImportProducts(ctx context.Context, req *v1.BulkImportProductsRequest) (*v1.BulkImportProductsResponse, error) {
	// In a real implementation, you would process the file content and import the products
	// For now, we'll just return a mock response
	return &v1.BulkImportProductsResponse{
		JobId:        uuid.New().String(),
		TotalRecords: 0,
		SuccessCount: 0,
		ErrorCount:   0,
		Status:       "pending",
	}, nil
}

// BulkExportProducts exports products to a file
func (s *ProductService) BulkExportProducts(ctx context.Context, req *v1.BulkExportProductsRequest) (*v1.BulkExportProductsResponse, error) {
	// In a real implementation, you would export the products to a file
	// For now, we'll just return a mock response
	return &v1.BulkExportProductsResponse{
		JobId:      uuid.New().String(),
		Status:     "pending",
		RecordCount: 0,
	}, nil
}

// RegisterService registers the ProductService with the gRPC server
func (s *ProductService) RegisterService(server interface{}) {
	if srv, ok := server.(v1.ProductServiceServer); ok {
		// This is a no-op since we're implementing the server interface directly
		_ = srv
	}
}
