package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hepstore/tribecart/proto/tribecart/v1"
	"google.golang.org/protobuf/types/known/timestamppb")

// Order represents an order in the database
type Order struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Status          string
	Subtotal        float64
	TaxAmount       float64
	ShippingCost    float64
	TotalAmount     float64
	PaymentMethod   string
	PaymentID       string
	ShippingAddress string
	TrackingNumber  string
	Carrier         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CancelledAt     *time.Time
	CancelledReason string
	Metadata        map[string]interface{}
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           uuid.UUID
	OrderID      uuid.UUID
	ProductID    uuid.UUID
	Quantity     int32
	UnitPrice    float64
	Subtotal     float64
	TaxAmount    float64
	TotalAmount  float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Metadata     map[string]interface{}
}

// OrderStatusHistory represents the history of status changes for an order
type OrderStatusHistory struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Status    string
	Notes     string
	CreatedAt time.Time
	Metadata  map[string]interface{}
}

// ToProto converts an Order to its protobuf representation
func (o *Order) ToProto() *tribecart_v1.Order {
	var cancelledAt *timestamppb.Timestamp
	if o.CancelledAt != nil {
		cancelledAt = timestamppb.New(*o.CancelledAt)
	}

	// Convert metadata to JSON string
	metadataJSON, _ := json.Marshal(o.Metadata)

	return &tribecart_v1.Order{
		Id:              o.ID.String(),
		UserId:          o.UserID.String(),
		Status:          tribecart_v1.OrderStatus(tribecart_v1.OrderStatus_value[o.Status]),
		Subtotal:        float32(o.Subtotal),
		TaxAmount:       float32(o.TaxAmount),
		ShippingCost:    float32(o.ShippingCost),
		TotalAmount:     float32(o.TotalAmount),
		PaymentMethod:   tribecart_v1.PaymentMethod(tribecart_v1.PaymentMethod_value[o.PaymentMethod]),
		PaymentId:       o.PaymentID,
		ShippingAddress: o.ShippingAddress,
		TrackingNumber:  o.TrackingNumber,
		Carrier:         o.Carrier,
		CreatedAt:       timestamppb.New(o.CreatedAt),
		UpdatedAt:       timestamppb.New(o.UpdatedAt),
		CancelledAt:     cancelledAt,
		CancelledReason: o.CancelledReason,
		Metadata:        string(metadataJSON),
	}
}

// OrderFromProto converts a protobuf Order to our database model
func OrderFromProto(pbOrder *tribecart_v1.Order) (*Order, error) {
	orderID, err := uuid.Parse(pbOrder.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %v", err)
	}

	userID, err := uuid.Parse(pbOrder.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var cancelledAt *time.Time
	if pbOrder.CancelledAt != nil {
		t := pbOrder.CancelledAt.AsTime()
		cancelledAt = &t
	}

	// Parse metadata
	var metadata map[string]interface{}
	if pbOrder.Metadata != "" {
		if err := json.Unmarshal([]byte(pbOrder.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata: %v", err)
		}
	}

	return &Order{
		ID:              orderID,
		UserID:          userID,
		Status:          pbOrder.Status.String(),
		Subtotal:        float64(pbOrder.Subtotal),
		TaxAmount:       float64(pbOrder.TaxAmount),
		ShippingCost:    float64(pbOrder.ShippingCost),
		TotalAmount:     float64(pbOrder.TotalAmount),
		PaymentMethod:   pbOrder.PaymentMethod.String(),
		PaymentID:       pbOrder.PaymentId,
		ShippingAddress: pbOrder.ShippingAddress,
		TrackingNumber:  pbOrder.TrackingNumber,
		Carrier:         pbOrder.Carrier,
		CreatedAt:       pbOrder.CreatedAt.AsTime(),
		UpdatedAt:       pbOrder.UpdatedAt.AsTime(),
		CancelledAt:     cancelledAt,
		CancelledReason: pbOrder.CancelledReason,
		Metadata:        metadata,
	}, nil
}

// OrderItemFromProto converts a protobuf OrderItem to our database model
func OrderItemFromProto(pbItem *tribecart_v1.OrderItem, orderID uuid.UUID) (*OrderItem, error) {
	itemID, err := uuid.Parse(pbItem.Id)
	if err != nil {
		itemID = uuid.New() // Generate new ID if not provided
	}

	productID, err := uuid.Parse(pbItem.ProductId)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Parse metadata
	var metadata map[string]interface{}
	if pbItem.Metadata != "" {
		if err := json.Unmarshal([]byte(pbItem.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("invalid item metadata: %v", err)
		}
	}

	return &OrderItem{
		ID:           itemID,
		OrderID:      orderID,
		ProductID:    productID,
		Quantity:     pbItem.Quantity,
		UnitPrice:    float64(pbItem.UnitPrice),
		Subtotal:     float64(pbItem.Subtotal),
		TaxAmount:    float64(pbItem.TaxAmount),
		TotalAmount:  float64(pbItem.TotalAmount),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     metadata,
	}, nil
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Order operations
	CreateOrder(ctx context.Context, order *Order, items []*OrderItem) error
	GetOrder(ctx context.Context, id uuid.UUID, includeItems, includeAddress bool) (*Order, []*OrderItem, error)
	ListOrders(ctx context.Context, userID *uuid.UUID, status *string, from, to *time.Time, limit, offset int) ([]*Order, int, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status, trackingNumber, carrier string) (*Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID, reason string, refundPayment bool) error
	
	// Order item operations
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*OrderItem, error)
	
	// Status history
	GetOrderStatusHistory(ctx context.Context, orderID uuid.UUID) ([]*OrderStatusHistory, error)
}

// DBOrderRepository implements OrderRepository using a SQL database
type DBOrderRepository struct {
	db *sql.DB
}

// NewDBOrderRepository creates a new DBOrderRepository
func NewDBOrderRepository(db *sql.DB) *DBOrderRepository {
	return &DBOrderRepository{db: db}
}

var _ OrderRepository = (*DBOrderRepository)(nil)

// CreateOrder creates a new order and its items in the database
func (r *DBOrderRepository) CreateOrder(ctx context.Context, order *Order, items []*OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	metadataJSON, _ := json.Marshal(order.Metadata)
	
	query := `
		INSERT INTO orders (
			id, user_id, status, subtotal, tax_amount, shipping_cost, total_amount,
			payment_method, payment_id, shipping_address_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.Status,
		order.Subtotal,
		order.TaxAmount,
		order.ShippingCost,
		order.TotalAmount,
		order.PaymentMethod,
		order.PaymentID,
		order.ShippingAddress,
		metadataJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (
			id, order_id, product_id, quantity, unit_price, 
			subtotal, tax_amount, total_amount, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, item := range items {
		itemMetadata, _ := json.Marshal(item.Metadata)
		
		_, err = tx.ExecContext(
			ctx,
			itemQuery,
			item.ID,
			order.ID, // Use the order's ID
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
			item.TaxAmount,
			item.TotalAmount,
			itemMetadata,
		)

		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Log the initial status
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO order_status_history (order_id, status, notes) VALUES ($1, $2, $3)`,
		order.ID,
		order.Status,
		"Order created",
	)

	if err != nil {
		return fmt.Errorf("failed to log order status: %w", err)
	}

	return tx.Commit()
}

// GetOrder retrieves an order by ID
func (r *DBOrderRepository) GetOrder(ctx context.Context, id uuid.UUID, includeItems, includeAddress bool) (*Order, []*OrderItem, error) {
	// Query order
	query := `
		SELECT 
			id, user_id, status, subtotal, tax_amount, shipping_cost, total_amount,
			payment_method, payment_id, shipping_address_id, tracking_number, carrier,
			created_at, updated_at, cancelled_at, cancelled_reason, metadata
		FROM orders
		WHERE id = $1
	`

	var order Order
	var metadataJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.ShippingCost,
		&order.TotalAmount,
		&order.PaymentMethod,
		&order.PaymentID,
		&order.ShippingAddress,
		&order.TrackingNumber,
		&order.Carrier,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CancelledAt,
		&order.CancelledReason,
		&metadataJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to query order: %w", err)
	}

	// Parse metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &order.Metadata); err != nil {
			return nil, nil, fmt.Errorf("failed to parse order metadata: %w", err)
		}
	}

	var items []*OrderItem
	if includeItems {
		items, err = r.GetOrderItems(ctx, id)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get order items: %w", err)
		}
	}

	return &order, items, nil
}

// GetOrderItems retrieves all items for an order
func (r *DBOrderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*OrderItem, error) {
	query := `
		SELECT 
			id, product_id, quantity, unit_price, subtotal, 
			tax_amount, total_amount, created_at, updated_at, metadata
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []*OrderItem
	for rows.Next() {
		var item OrderItem
		var metadataJSON []byte

		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.TaxAmount,
			&item.TotalAmount,
			&item.CreatedAt,
			&item.UpdatedAt,
			&metadataJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		// Parse metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
				return nil, fmt.Errorf("failed to parse item metadata: %w", err)
			}
		}

		item.OrderID = orderID
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return items, nil
}

// ListOrders retrieves a paginated list of orders with optional filtering
func (r *DBOrderRepository) ListOrders(
	ctx context.Context,
	userID *uuid.UUID,
	status *string,
	from, to *time.Time,
	limit, offset int,
) ([]*Order, int, error) {
	// Build the base query and count query
	baseQuery := `
		SELECT 
			o.id, o.user_id, o.status, o.subtotal, o.tax_amount, o.shipping_cost, 
			o.total_amount, o.payment_method, o.payment_id, o.shipping_address_id,
			o.tracking_number, o.carrier, o.created_at, o.updated_at, 
			o.cancelled_at, o.cancelled_reason, o.metadata
		FROM orders o
	`

	countQuery := `SELECT COUNT(*) FROM orders o`

	// Build WHERE clause
	var whereClause string
	var args []interface{}
	argNum := 1

	if userID != nil {
		whereClause = fmt.Sprintf("WHERE o.user_id = $%d", argNum)
		args = append(args, *userID)
		argNum++
	}

	if status != nil {
		if whereClause == "" {
			whereClause = "WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("o.status = $%d", argNum)
		args = append(args, *status)
		argNum++
	}

	if from != nil {
		if whereClause == "" {
			whereClause = "WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("o.created_at >= $%d", argNum)
		args = append(args, *from)
		argNum++
	}

	if to != nil {
		if whereClause == "" {
			whereClause = "WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("o.created_at <= $%d", argNum)
		args = append(args, *to)
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Add pagination to the base query
	orderClause := " ORDER BY o.created_at DESC"
	pagination := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, baseQuery+whereClause+orderClause+pagination, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		var metadataJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.Subtotal,
			&order.TaxAmount,
			&order.ShippingCost,
			&order.TotalAmount,
			&order.PaymentMethod,
			&order.PaymentID,
			&order.ShippingAddress,
			&order.TrackingNumber,
			&order.Carrier,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CancelledAt,
			&order.CancelledReason,
			&metadataJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}

		// Parse metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &order.Metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to parse order metadata: %w", err)
			}
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates an order's status and optionally tracking information
func (r *DBOrderRepository) UpdateOrderStatus(
	ctx context.Context,
	id uuid.UUID,
	status, trackingNumber, carrier string,
) (*Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update the order
	query := `
		UPDATE orders
		SET 
			status = $1,
			tracking_number = COALESCE($2, tracking_number),
			carrier = COALESCE($3, carrier),
			updated_at = NOW()
		WHERE id = $4
		RETURNING 
			id, user_id, status, subtotal, tax_amount, shipping_cost, 
			total_amount, payment_method, payment_id, shipping_address_id,
			tracking_number, carrier, created_at, updated_at, 
			cancelled_at, cancelled_reason, metadata
	`

	var order Order
	var metadataJSON []byte

	err = tx.QueryRowContext(
		ctx,
		query,
		status,
		trackingNumber,
		carrier,
		id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.ShippingCost,
		&order.TotalAmount,
		&order.PaymentMethod,
		&order.PaymentID,
		&order.ShippingAddress,
		&order.TrackingNumber,
		&order.Carrier,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CancelledAt,
		&order.CancelledReason,
		&metadataJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Parse metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &order.Metadata); err != nil {
			return nil, fmt.Errorf("failed to parse order metadata: %w", err)
		}
	}

	// Log the status change
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO order_status_history (order_id, status, notes) VALUES ($1, $2, $3)`,
		id,
		status,
		"Order status updated",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to log status change: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

// CancelOrder cancels an order and optionally processes a refund
func (r *DBOrderRepository) CancelOrder(
	ctx context.Context,
	id uuid.UUID,
	reason string,
	refundPayment bool,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, get the current status to validate the transition
	var currentStatus string
	err = tx.QueryRowContext(
		ctx,
		`SELECT status FROM orders WHERE id = $1 FOR UPDATE`,
		id,
	).Scan(&currentStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order not found: %w", err)
		}
		return fmt.Errorf("failed to get order status: %w", err)
	}

	// Validate that the order can be cancelled
	if currentStatus == "CANCELLED" || currentStatus == "REFUNDED" {
		return fmt.Errorf("order is already %s", currentStatus)
	}

	// Update the order status to CANCELLED
	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE orders
			SET 
				status = 'CANCELLED',
				cancelled_at = NOW(),
				cancelled_reason = $1,
				updated_at = NOW()
			WHERE id = $2
		`,
		reason,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Log the cancellation
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO order_status_history 
			(order_id, status, notes) 
			VALUES ($1, 'CANCELLED', $2)
		`,
		id,
		reason,
	)

	if err != nil {
		return fmt.Errorf("failed to log cancellation: %w", err)
	}

	// TODO: Process refund if requested
	// This would integrate with the payment service
	if refundPayment {
		// Implement refund logic here
	}

	return tx.Commit()
}

// GetOrderStatusHistory retrieves the status history for an order
func (r *DBOrderRepository) GetOrderStatusHistory(
	ctx context.Context,
	orderID uuid.UUID,
) ([]*OrderStatusHistory, error) {
	query := `
		SELECT 
			id, status, notes, created_at, metadata
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query status history: %w", err)
	}
	defer rows.Close()

	var history []*OrderStatusHistory
	for rows.Next() {
		var h OrderStatusHistory
		var metadataJSON []byte

		err := rows.Scan(
			&h.ID,
			&h.Status,
			&h.Notes,
			&h.CreatedAt,
			&metadataJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan status history: %w", err)
		}

		// Parse metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &h.Metadata); err != nil {
				return nil, fmt.Errorf("failed to parse status metadata: %w", err)
			}
		}

		h.OrderID = orderID
		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status history: %w", err)
	}

	return history, nil
}
