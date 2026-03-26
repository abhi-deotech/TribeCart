package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hepstore/tribecart/proto/tribecart/v1"
	"github.com/hepstore/tribecart/services/orders/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb")

// OrderService implements the gRPC OrderService server
type OrderService struct {
	repo models.OrderRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(repo models.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(
	ctx context.Context,
	req *tribecart_v1.CreateOrderRequest,
) (*tribecart_v1.Order, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one order item is required")
	}

	if req.ShippingAddressId == "" {
		return nil, status.Error(codes.InvalidArgument, "shipping_address_id is required")
	}

	// Convert protobuf items to model items
	var orderItems []*models.OrderItem
	orderID := uuid.New()

	for _, item := range req.Items {
		orderItem, err := models.OrderItemFromProto(item, orderID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid order item: %v", err)
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order model
	order := &models.Order{
		ID:              orderID,
		UserID:          uuid.MustParse(req.UserId),
		Status:          tribecart_v1.OrderStatus_ORDER_STATUS_PENDING.String(),
		Subtotal:        float64(req.Subtotal),
		TaxAmount:       float64(req.TaxAmount),
		ShippingCost:    float64(req.ShippingCost),
		TotalAmount:     float64(req.TotalAmount),
		PaymentMethod:   req.PaymentMethod.String(),
		PaymentID:       req.PaymentId,
		ShippingAddress: req.ShippingAddressId,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Add metadata if provided
	if req.Metadata != nil {
		metadata, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid metadata: %v", err)
		}
		order.Metadata = make(map[string]interface{})
		if err := json.Unmarshal(metadata, &order.Metadata); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid metadata format: %v", err)
		}
	}

	// Save to database
	if err := s.repo.CreateOrder(ctx, order, orderItems); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Convert back to protobuf and return
	return order.ToProto(), nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(
	ctx context.Context,
	req *tribecart_v1.GetOrderRequest,
) (*tribecart_v1.Order, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	// Get order from repository
	order, items, err := s.repo.GetOrder(ctx, orderID, req.IncludeItems, req.IncludeAddress)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	// Convert to protobuf
	pbOrder := order.ToProto()

	// Add items if requested
	if req.IncludeItems && len(items) > 0 {
		for _, item := range items {
			pbOrder.Items = append(pbOrder.Items, &tribecart_v1.OrderItem{
				Id:          item.ID.String(),
				ProductId:   item.ProductID.String(),
				Quantity:    item.Quantity,
				UnitPrice:   float32(item.UnitPrice),
				Subtotal:    float32(item.Subtotal),
				TaxAmount:   float32(item.TaxAmount),
				TotalAmount: float32(item.TotalAmount),
				Metadata:    string(mustMarshalJSON(item.Metadata)),
			})
		}
	}

	return pbOrder, nil
}

// ListOrders retrieves a list of orders with optional filtering
func (s *OrderService) ListOrders(
	ctx context.Context,
	req *tribecart_v1.ListOrdersRequest,
) (*tribecart_v1.ListOrdersResponse, error) {
	// Parse pagination
	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50 // Default page size
	}

	offset := int(req.Page) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Parse filters
	var userID *uuid.UUID
	if req.UserId != "" {
		parsedUserID, err := uuid.Parse(req.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
		}
		userID = &parsedUserID
	}

	var statusFilter *string
	if req.Status != tribecart_v1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		statusStr := req.Status.String()
		statusFilter = &statusStr
	}

	var fromDate, toDate *time.Time
	if req.FromDate != nil {
		t := req.FromDate.AsTime()
		fromDate = &t
	}
	if req.ToDate != nil {
		t := req.ToDate.AsTime()
		toDate = &t
	}

	// Get orders from repository
	orders, total, err := s.repo.ListOrders(ctx, userID, statusFilter, fromDate, toDate, pageSize, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	// Convert to protobuf
	var pbOrders []*tribecart_v1.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, order.ToProto())
	}

	// Calculate next page token if there are more results
	var nextPageToken string
	if (offset + len(orders)) < total {
		nextPageToken = fmt.Sprintf("%d", req.Page+1)
	}

	return &tribecart_v1.ListOrdersResponse{
		Orders:         pbOrders,
		NextPageToken: nextPageToken,
		TotalSize:     int32(total),
	}, nil
}

// UpdateOrderStatus updates an order's status
func (s *OrderService) UpdateOrderStatus(
	ctx context.Context,
	req *tribecart_v1.UpdateOrderStatusRequest,
) (*tribecart_v1.Order, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	// Validate status
	statusStr := req.Status.String()
	if _, ok := tribecart_v1.OrderStatus_value[statusStr]; !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid order status")
	}

	// Update order status in repository
	order, err := s.repo.UpdateOrderStatus(
		ctx,
		orderID,
		statusStr,
		req.TrackingNumber,
		req.Carrier,
	)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return order.ToProto(), nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(
	ctx context.Context,
	req *tribecart_v1.CancelOrderRequest,
) (*tribecart_v1.Order, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	orderID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	// Cancel order in repository
	if err := s.repo.CancelOrder(ctx, orderID, req.Reason, req.RefundPayment); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		if errors.Is(err, models.ErrInvalidStatus) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	// Return the updated order
	order, _, err := s.repo.GetOrder(ctx, orderID, false, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated order: %v", err)
	}

	return order.ToProto(), nil
}

// GetOrderHistory retrieves the status history for an order
func (s *OrderService) GetOrderHistory(
	ctx context.Context,
	req *tribecart_v1.GetOrderHistoryRequest,
) (*tribecart_v1.OrderHistoryResponse, error) {
	// Validate request
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
	}

	// Get status history from repository
	history, err := s.repo.GetOrderStatusHistory(ctx, orderID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order history: %v", err)
	}

	// Convert to protobuf
	var pbHistory []*tribecart_v1.OrderStatusUpdate
	for _, h := range history {
		pbHistory = append(pbHistory, &tribecart_v1.OrderStatusUpdate{
			Id:        h.ID.String(),
			OrderId:   h.OrderID.String(),
			Status:    tribecart_v1.OrderStatus(tribecart_v1.OrderStatus_value[h.Status]),
			Notes:     h.Notes,
			CreatedAt: timestamppb.New(h.CreatedAt),
			Metadata:  string(mustMarshalJSON(h.Metadata)),
		})
	}

	return &tribecart_v1.OrderHistoryResponse{
		Updates: pbHistory,
	}, nil
}

// Helper function to marshal JSON, panics on error (should never happen with valid input)
func mustMarshalJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
