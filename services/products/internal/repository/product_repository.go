package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tribecart/proto/tribecart/v1"
)

var ErrNotFound = errors.New("not found")

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	// CreateProduct creates a new product in the database
	CreateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error)

	// GetProduct retrieves a product by ID
	GetProduct(ctx context.Context, id string, includeDeleted bool) (*pb.Product, error)

	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, product *pb.Product, updateMask []string) (*pb.Product, error)

	// DeleteProduct soft deletes a product by ID
	DeleteProduct(ctx context.Context, id string, force bool) error

	// ListProducts retrieves a paginated list of products with filtering
	ListProducts(ctx context.Context, req *pb.ListProductsRequest) ([]*pb.Product, int32, error)

	// UpdateStock updates the stock quantity for a product or variant
	UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.Product, error)

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CommitTx commits a transaction
	CommitTx(tx *sql.Tx) error

	// RollbackTx rolls back a transaction
	RollbackTx(tx *sql.Tx) error
}

// PostgresProductRepository implements ProductRepository for PostgreSQL
type PostgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository creates a new PostgresProductRepository
func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// scanProduct scans a product from a sql.Row or sql.Rows
func scanProduct(row interface{ Scan(...interface{}) error }) (*pb.Product, error) {
	var (
		id, sellerID, name, description, sku, barcode, taxClassID, seoTitle, seoDescription, seoKeywords sql.NullString
		price, salePrice, costPrice, weight, length, width, height                                                                  sql.NullFloat64
		stockQuantity, minStockLevel                                                                                                sql.NullInt32
		trackInventory, isFeatured, isVisible, requiresShipping, isTaxable                                                           bool
		status                                                                                                                                  int32
		productType                                                                                                                                  int32
		createdAt, updatedAt, deletedAt                                                                                           sql.NullTime
	)

	err := row.Scan(
		&id, &sellerID, &name, &description, &sku, &barcode,
		&price, &salePrice, &costPrice, &stockQuantity, &trackInventory,
		&minStockLevel, &weight, &length, &width, &height, &productType,
		&status, &isFeatured, &isVisible, &requiresShipping, &isTaxable,
		&taxClassID, &seoTitle, &seoDescription, &seoKeywords,
		&createdAt, &updatedAt, &deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	product := &pb.Product{
		Id:              id.String,
		SellerId:        sellerID.String,
		Name:            name.String,
		Description:     description.String,
		Sku:             sku.String,
		Barcode:         barcode.String,
		Price:           price.Float64,
		Status:          pb.ProductStatus(status),
		Type:            pb.ProductType(productType),
		IsFeatured:      isFeatured,
		IsVisible:       isVisible,
		RequiresShipping: requiresShipping,
		IsTaxable:       isTaxable,
	}

	if salePrice.Valid {
		product.SalePrice = salePrice.Float64
	}

	if costPrice.Valid {
		product.CostPrice = costPrice.Float64
	}

	if stockQuantity.Valid {
		product.StockQuantity = stockQuantity.Int32
	}

	product.TrackInventory = trackInventory

	if minStockLevel.Valid {
		product.MinStockLevel = minStockLevel.Int32
	}

	if weight.Valid {
		product.Weight = weight.Float64
	}

	if length.Valid {
		product.Length = length.Float64
	}

	if width.Valid {
		product.Width = width.Float64
	}

	if height.Valid {
		product.Height = height.Float64
	}

	if taxClassID.Valid {
		product.TaxClassId = taxClassID.String
	}

	if seoTitle.Valid {
		product.SeoTitle = seoTitle.String
	}

	if seoDescription.Valid {
		product.SeoDescription = seoDescription.String
	}

	if seoKeywords.Valid {
		product.SeoKeywords = seoKeywords.String
	}

	if createdAt.Valid {
		product.CreatedAt = timestamppb.New(createdAt.Time)
	}

	if updatedAt.Valid {
		product.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	if deletedAt.Valid {
		product.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return product, nil
}

// CreateProduct creates a new product in the database
func (r *PostgresProductRepository) CreateProduct(ctx context.Context, product *pb.Product) (*pb.Product, error) {
	if product.Id == "" {
		product.Id = uuid.New().String()
	}

	query := `
		INSERT INTO products (
			id, seller_id, name, description, sku, barcode, price, sale_price, 
			cost_price, stock_quantity, track_inventory, min_stock_level, 
			weight, length, width, height, type, status, is_featured, 
			is_visible, requires_shipping, is_taxable, tax_class_id, 
			seo_title, seo_description, seo_keywords
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
			$11, $12, $13, $14, $15, $16, $17, $18, $19, 
			$20, $21, $22, $23, $24, $25, $26
		) RETURNING *
	`

	row := r.db.QueryRowContext(ctx, query,
		product.Id,
		sql.NullString{String: product.SellerId, Valid: product.SellerId != ""},
		sql.NullString{String: product.Name, Valid: product.Name != ""},
		sql.NullString{String: product.Description, Valid: product.Description != ""},
		sql.NullString{String: product.Sku, Valid: product.Sku != ""},
		sql.NullString{String: product.Barcode, Valid: product.Barcode != ""},
		sql.NullFloat64{Float64: product.Price, Valid: true},
		sql.NullFloat64{Float64: product.SalePrice, Valid: product.SalePrice != 0},
		sql.NullFloat64{Float64: product.CostPrice, Valid: product.CostPrice != 0},
		sql.NullInt32{Int32: product.StockQuantity, Valid: true},
		product.TrackInventory,
		sql.NullInt32{Int32: product.MinStockLevel, Valid: product.MinStockLevel > 0},
		sql.NullFloat64{Float64: product.Weight, Valid: product.Weight > 0},
		sql.NullFloat64{Float64: product.Length, Valid: product.Length > 0},
		sql.NullFloat64{Float64: product.Width, Valid: product.Width > 0},
		sql.NullFloat64{Float64: product.Height, Valid: product.Height > 0},
		product.Type,
		product.Status,
		product.IsFeatured,
		product.IsVisible,
		product.RequiresShipping,
		product.IsTaxable,
		sql.NullString{String: product.TaxClassId, Valid: product.TaxClassId != ""},
		sql.NullString{String: product.SeoTitle, Valid: product.SeoTitle != ""},
		sql.NullString{String: product.SeoDescription, Valid: product.SeoDescription != ""},
		sql.NullString{String: product.SeoKeywords, Valid: product.SeoKeywords != ""},
	)

	return scanProduct(row)
}

// GetProduct retrieves a product by ID
func (r *PostgresProductRepository) GetProduct(ctx context.Context, id string, includeDeleted bool) (*pb.Product, error) {
	query := `
		SELECT * FROM products 
		WHERE id = $1
	`

	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	row := r.db.QueryRowContext(ctx, query, id)
	return scanProduct(row)
}

// UpdateProduct updates an existing product
func (r *PostgresProductRepository) UpdateProduct(ctx context.Context, product *pb.Product, updateMask []string) (*pb.Product, error) {
	// This is a simplified implementation. In a real-world scenario, you would
	// use the updateMask to build a dynamic SQL query that only updates the specified fields.
	// For simplicity, we'll update all fields here.

	query := `
		UPDATE products SET
			seller_id = $2,
			name = $3,
			description = $4,
			sku = $5,
			barcode = $6,
			price = $7,
			sale_price = $8,
			cost_price = $9,
			stock_quantity = $10,
			track_inventory = $11,
			min_stock_level = $12,
			weight = $13,
			length = $14,
			width = $15,
			height = $16,
			type = $17,
			status = $18,
			is_featured = $19,
			is_visible = $20,
			requires_shipping = $21,
			is_taxable = $22,
			tax_class_id = $23,
			seo_title = $24,
			seo_description = $25,
			seo_keywords = $26,
			updated_at = NOW()
		WHERE id = $1
		RETURNING *
	`

	row := r.db.QueryRowContext(ctx, query,
		product.Id,
		sql.NullString{String: product.SellerId, Valid: product.SellerId != ""},
		sql.NullString{String: product.Name, Valid: product.Name != ""},
		sql.NullString{String: product.Description, Valid: product.Description != ""},
		sql.NullString{String: product.Sku, Valid: product.Sku != ""},
		sql.NullString{String: product.Barcode, Valid: product.Barcode != ""},
		sql.NullFloat64{Float64: product.Price, Valid: true},
		sql.NullFloat64{Float64: product.SalePrice, Valid: product.SalePrice != 0},
		sql.NullFloat64{Float64: product.CostPrice, Valid: product.CostPrice != 0},
		sql.NullInt32{Int32: product.StockQuantity, Valid: true},
		product.TrackInventory,
		sql.NullInt32{Int32: product.MinStockLevel, Valid: product.MinStockLevel > 0},
		sql.NullFloat64{Float64: product.Weight, Valid: product.Weight > 0},
		sql.NullFloat64{Float64: product.Length, Valid: product.Length > 0},
		sql.NullFloat64{Float64: product.Width, Valid: product.Width > 0},
		sql.NullFloat64{Float64: product.Height, Valid: product.Height > 0},
		product.Type,
		product.Status,
		product.IsFeatured,
		product.IsVisible,
		product.RequiresShipping,
		product.IsTaxable,
		sql.NullString{String: product.TaxClassId, Valid: product.TaxClassId != ""},
		sql.NullString{String: product.SeoTitle, Valid: product.SeoTitle != ""},
		sql.NullString{String: product.SeoDescription, Valid: product.SeoDescription != ""},
		sql.NullString{String: product.SeoKeywords, Valid: product.SeoKeywords != ""},
	)

	return scanProduct(row)
}

// DeleteProduct soft deletes a product by ID
func (r *PostgresProductRepository) DeleteProduct(ctx context.Context, id string, force bool) error {
	if force {
		// Hard delete
		_, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
		return err
	}

	// Soft delete
	_, err := r.db.ExecContext(ctx, "UPDATE products SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ListProducts retrieves a paginated list of products with filtering
func (r *PostgresProductRepository) ListProducts(ctx context.Context, req *pb.ListProductsRequest) ([]*pb.Product, int32, error) {
	// Build the base query
	query := "SELECT * FROM products WHERE 1=1"
	var args []interface{}
	argPos := 1

	// Apply filters
	if req.SellerId != "" {
		query += fmt.Sprintf(" AND seller_id = $%d", argPos)
		args = append(args, req.SellerId)
		argPos++
	}

	if len(req.CategoryIds) > 0 {
		// Assuming a product_categories table exists for many-to-many relationship
		query += fmt.Sprintf(" AND id IN (SELECT product_id FROM product_categories WHERE category_id = ANY($%d))", argPos)
		args = append(args, req.CategoryIds)
		argPos++
	}

	if req.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argPos)
		args = append(args, req.MinPrice)
		argPos++
	}

	if req.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argPos)
		args = append(args, req.MaxPrice)
		argPos++
	}

	if req.Status != pb.ProductStatus_PRODUCT_STATUS_UNSPECIFIED {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, req.Status)
		argPos++
	}

	if req.InStock {
		query += " AND (stock_quantity > 0 OR track_inventory = false)"
	}

	if req.Featured {
		query += " AND is_featured = true"
	}

	if req.SearchQuery != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR sku = $%d)", argPos, argPos, argPos)
		searchTerm := "%" + req.SearchQuery + "%"
		args = append(args, searchTerm, searchTerm, req.SearchQuery)
		argPos += 3
	}

	if !req.IncludeDeleted {
		query += " AND deleted_at IS NULL"
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_query"
	var totalCount int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if req.SortBy != "" {
		sortField := strings.TrimPrefix(req.SortBy, "-")
		sortOrder := "ASC"
		if strings.HasPrefix(req.SortBy, "-") {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}

	// Apply pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan the results
	var products []*pb.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

// UpdateStock updates the stock quantity for a product or variant
func (r *PostgresProductRepository) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.Product, error) {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			r.RollbackTx(tx)
		}
	}()

	// Update product stock
	query := `
		UPDATE products 
		SET stock_quantity = stock_quantity + $1,
		    updated_at = NOW()
		WHERE id = $2
		RETURNING *
	`

	row := tx.QueryRowContext(ctx, query, req.Quantity, req.ProductId)
	product, err := scanProduct(row)
	if err != nil {
		return nil, err
	}

	// Update variant stock if variant_id is provided
	if req.VariantId != "" {
		variantQuery := `
			UPDATE product_variants
			SET stock_quantity = stock_quantity + $1,
			    updated_at = NOW()
			WHERE id = $2 AND product_id = $3
			RETURNING *
		`
		// In a real implementation, you would scan the variant here
		_, err = tx.ExecContext(ctx, variantQuery, req.Quantity, req.VariantId, req.ProductId)
		if err != nil {
			return nil, err
		}

		// Update the product's stock quantity to reflect the sum of all variants
		if err = r.updateProductStockFromVariants(ctx, tx, req.ProductId); err != nil {
			return nil, err
		}

		// Refetch the product to get updated stock quantity
		row = tx.QueryRowContext(ctx, "SELECT * FROM products WHERE id = $1", req.ProductId)
		product, err = scanProduct(row)
		if err != nil {
			return nil, err
		}
	}

	// Record stock movement
	err = r.recordStockMovement(ctx, tx, &StockMovement{
		ProductID:   req.ProductId,
		VariantID:   req.VariantId,
		Quantity:    req.Quantity,
		ReferenceID: req.ReferenceId,
		Reason:      req.Reason,
	})
	if err != nil {
		return nil, err
	}

	if err = r.CommitTx(tx); err != nil {
		return nil, err
	}

	return product, nil
}

// updateProductStockFromVariants updates the product's stock quantity based on the sum of its variants
func (r *PostgresProductRepository) updateProductStockFromVariants(ctx context.Context, tx *sql.Tx, productID string) error {
	query := `
		UPDATE products p
		SET stock_quantity = (
			SELECT COALESCE(SUM(stock_quantity), 0)
			FROM product_variants
			WHERE product_id = $1
		)
		WHERE id = $1
	`
	_, err := tx.ExecContext(ctx, query, productID)
	return err
}

// StockMovement represents a stock movement record
type StockMovement struct {
	ID          string
	ProductID   string
	VariantID   string
	Quantity    int32
	ReferenceID string
	Reason      string
	CreatedAt   *timestamppb.Timestamp
}

// recordStockMovement records a stock movement in the database
func (r *PostgresProductRepository) recordStockMovement(ctx context.Context, tx *sql.Tx, movement *StockMovement) error {
	query := `
		INSERT INTO stock_movements (
			id, product_id, variant_id, quantity, reference_id, reason
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := tx.ExecContext(ctx, query,
		uuid.New().String(),
		movement.ProductID,
		sql.NullString{String: movement.VariantID, Valid: movement.VariantID != ""},
		movement.Quantity,
		sql.NullString{String: movement.ReferenceID, Valid: movement.ReferenceID != ""},
		movement.Reason,
	)

	return err
}

// BeginTx starts a new transaction
func (r *PostgresProductRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// CommitTx commits a transaction
func (r *PostgresProductRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

// RollbackTx rolls back a transaction
func (r *PostgresProductRepository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}
