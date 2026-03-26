package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hepstore/tribecart/proto/tribecart/v1"
)

func TestCreateOrder(t *testing.T) {
	s := NewServer()

	tests := []struct {
		name    string
		req     *pb.CreateOrderRequest
		wantErr bool
		errCode  codes.Code
	}{
		{
			name: "valid order",
			req: &pb.CreateOrderRequest{
				UserId: "user123",
				Items: []*pb.OrderItem{
					{ProductId: "prod1", Name: "Test Product", Quantity: 2, UnitPrice: 10.0, TotalPrice: 20.0},
				},
				ShippingAddressId: "addr1",
				PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
				PaymentId:        "pay123",
			},
			wantErr: false,
		},
		{
			name: "empty items",
			req: &pb.CreateOrderRequest{
				UserId: "user123",
				Items:  []*pb.OrderItem{},
			},
			wantErr: true,
			errCode:  codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.CreateOrder(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode != codes.OK {
					s, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, s.Code())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Id)
			assert.Equal(t, tt.req.UserId, resp.UserId)
			assert.Equal(t, pb.OrderStatus_ORDER_STATUS_PENDING, resp.Status)
			assert.NotNil(t, resp.OrderedAt)
		})
	}
}

func TestGetOrder(t *testing.T) {
	s := NewServer()

	// Create a test order
	createReq := &pb.CreateOrderRequest{
		UserId: "user123",
		Items: []*pb.OrderItem{
			{ProductId: "prod1", Name: "Test Product", Quantity: 1, UnitPrice: 10.0, TotalPrice: 10.0},
		},
		ShippingAddressId: "addr1",
		PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}

	createdOrder, err := s.CreateOrder(context.Background(), createReq)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		orderID string
		wantErr bool
		errCode  codes.Code
	}{
		{
			name:    "existing order",
			orderID: createdOrder.Id,
			wantErr: false,
		},
		{
			name:    "non-existent order",
			orderID: "nonexistent",
			wantErr: true,
			errCode:  codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.GetOrderRequest{Id: tt.orderID}
			resp, err := s.GetOrder(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode != codes.OK {
					s, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, s.Code())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.orderID, resp.Id)
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	s := NewServer()

	// Create a test order
	createReq := &pb.CreateOrderRequest{
		UserId: "user123",
		Items: []*pb.OrderItem{
			{ProductId: "prod1", Name: "Test Product", Quantity: 1, UnitPrice: 10.0, TotalPrice: 10.0},
		},
		ShippingAddressId: "addr1",
		PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}

	createdOrder, err := s.CreateOrder(context.Background(), createReq)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		orderID    string
		status     pb.OrderStatus
		tracking   string
		carrier    string
		expected   pb.OrderStatus
		wantErr    bool
		errCode    codes.Code
		errMessage string
	}{
		{
			name:     "valid status update to processing",
			orderID:  createdOrder.Id,
			status:   pb.OrderStatus_ORDER_STATUS_PROCESSING,
			expected: pb.OrderStatus_ORDER_STATUS_PROCESSING,
			wantErr:  false,
		},
		{
			name:     "valid status update to shipped",
			orderID:  createdOrder.Id,
			status:   pb.OrderStatus_ORDER_STATUS_SHIPPED,
			tracking: "TRACK123",
			carrier:  "FEDEX",
			expected: pb.OrderStatus_ORDER_STATUS_SHIPPED,
			wantErr:  false,
		},
		{
			name:       "invalid status transition",
			orderID:    createdOrder.Id,
			status:     pb.OrderStatus_ORDER_STATUS_DELIVERED,
			expected:   pb.OrderStatus_ORDER_STATUS_PENDING,
			wantErr:    true,
			errCode:    codes.InvalidArgument,
			errMessage: "invalid status transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.UpdateOrderStatusRequest{
				Id:            tt.orderID,
				Status:        tt.status,
				TrackingNumber: tt.tracking,
				Carrier:       tt.carrier,
			}

			resp, err := s.UpdateOrderStatus(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode != codes.OK {
					s, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, s.Code())
					if tt.errMessage != "" {
						assert.Contains(t, s.Message(), tt.errMessage)
					}
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, resp.Status)

			// Verify tracking info if provided
			if tt.tracking != "" {
				assert.Equal(t, tt.tracking, resp.TrackingNumber)
				assert.Equal(t, tt.carrier, resp.Carrier)
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	s := NewServer()

	// Create a test order
	createReq := &pb.CreateOrderRequest{
		UserId: "user123",
		Items: []*pb.OrderItem{
			{ProductId: "prod1", Name: "Test Product", Quantity: 1, UnitPrice: 10.0, TotalPrice: 10.0},
		},
		ShippingAddressId: "addr1",
		PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	}

	createdOrder, err := s.CreateOrder(context.Background(), createReq)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		orderID    string
		reason     string
		wantErr    bool
		errCode    codes.Code
		errMessage string
	}{
		{
			name:    "valid cancellation",
			orderID: createdOrder.Id,
			reason:  "Changed my mind",
			wantErr: false,
		},
		{
			name:       "non-existent order",
			orderID:    "nonexistent",
			reason:     "Test",
			wantErr:    true,
			errCode:    codes.NotFound,
			errMessage: "order not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CancelOrderRequest{
				Id:     tt.orderID,
				Reason: tt.reason,
			}

			_, err := s.CancelOrder(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errCode != codes.OK {
					s, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, tt.errCode, s.Code())
					if tt.errMessage != "" {
						assert.Contains(t, s.Message(), tt.errMessage)
					}
				}
				return
			}

			assert.NoError(t, err)

			// Verify the order was cancelled
			getReq := &pb.GetOrderRequest{Id: tt.orderID}
			resp, err := s.GetOrder(context.Background(), getReq)
			assert.NoError(t, err)
			assert.Equal(t, pb.OrderStatus_ORDER_STATUS_CANCELLED, resp.Status)
			assert.Equal(t, tt.reason, resp.CancellationReason)
			assert.NotNil(t, resp.CancelledAt)
		})
	}
}

func TestListOrders(t *testing.T) {
	s := NewServer()

	// Create test data
	user1 := "user1"
	user2 := "user2"
	now := time.Now()

	// Create orders for user1
	for i := 0; i < 3; i++ {
		req := &pb.CreateOrderRequest{
			UserId: user1,
			Items: []*pb.OrderItem{
				{ProductId: "prod1", Name: "Product 1", Quantity: 1, UnitPrice: 10.0, TotalPrice: 10.0},
			},
			ShippingAddressId: "addr1",
			PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}
		_, err := s.CreateOrder(context.Background(), req)
		assert.NoError(t, err)
	}

	// Create orders for user2
	for i := 0; i < 2; i++ {
		req := &pb.CreateOrderRequest{
			UserId: user2,
			Items: []*pb.OrderItem{
				{ProductId: "prod2", Name: "Product 2", Quantity: 2, UnitPrice: 15.0, TotalPrice: 30.0},
			},
			ShippingAddressId: "addr2",
			PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_UPI,
		}
		_, err := s.CreateOrder(context.Background(), req)
		assert.NoError(t, err)

		// Add a small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	tests := []struct {
		name      string
		userID    string
		status    pb.OrderStatus
		expCount  int
		wantError bool
	}{
		{
			name:     "all orders",
			expCount: 5,
		},
		{
			name:     "user1 orders",
			userID:   user1,
			expCount: 3,
		},
		{
			name:     "user2 orders",
			userID:   user2,
			expCount: 2,
		},
		{
			name:     "non-existent user",
			userID:   "nonexistent",
			expCount: 0,
		},
		{
			name:     "filter by status",
			status:   pb.OrderStatus_ORDER_STATUS_PENDING,
			expCount: 5, // All orders should be pending
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ListOrdersRequest{
				UserId: tt.userID,
				Status: tt.status,
			}

			resp, err := s.ListOrders(context.Background(), req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expCount, len(resp.Orders))

			// Verify the orders match the filters
			for _, order := range resp.Orders {
				if tt.userID != "" {
					assert.Equal(t, tt.userID, order.UserId)
				}
				if tt.status != pb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
					assert.Equal(t, tt.status, order.Status)
				}
			}
		})
	}
}

func TestGetOrderHistory(t *testing.T) {
	s := NewServer()

	// Create test data
	userID := "test_user"
	now := time.Now()

	// Create orders with different statuses
	statuses := []pb.OrderStatus{
		pb.OrderStatus_ORDER_STATUS_PENDING,
		pb.OrderStatus_ORDER_STATUS_PROCESSING,
		pb.OrderStatus_ORDER_STATUS_SHIPPED,
		pb.OrderStatus_ORDER_STATUS_DELIVERED,
		pb.OrderStatus_ORDER_STATUS_CANCELLED,
	}

	for i, status := range statuses {
		req := &pb.CreateOrderRequest{
			UserId: userID,
			Items: []*pb.OrderItem{
				{ProductId: "prod1", Name: "Product 1", Quantity: 1, UnitPrice: 10.0, TotalPrice: 10.0},
			},
			ShippingAddressId: "addr1",
			PaymentMethod:    pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		}

		order, err := s.CreateOrder(context.Background(), req)
		assert.NoError(t, err)

		// Update order status if not pending
		if status != pb.OrderStatus_ORDER_STATUS_PENDING {
			_, err = s.UpdateOrderStatus(context.Background(), &pb.UpdateOrderStatusRequest{
				Id:     order.Id,
				Status: status,
			})
			assert.NoError(t, err)
		}

		// Add a small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	tests := []struct {
		name     string
		userID   string
		status   pb.OrderStatus
		expCount int
	}{
		{
			name:     "all orders",
			userID:   userID,
			expCount: len(statuses),
		},
		{
			name:     "filter by status - pending",
			userID:   userID,
			status:   pb.OrderStatus_ORDER_STATUS_PENDING,
			expCount: 1,
		},
		{
			name:     "filter by status - delivered",
			userID:   userID,
			status:   pb.OrderStatus_ORDER_STATUS_DELIVERED,
			expCount: 1,
		},
		{
			name:     "non-existent user",
			userID:   "nonexistent",
			expCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.GetOrderHistoryRequest{
				UserId: tt.userID,
				Status: tt.status,
			}

			resp, err := s.GetOrderHistory(context.Background(), req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expCount, len(resp.Orders))

			// Verify the orders match the filters
			for _, order := range resp.Orders {
				assert.Equal(t, tt.userID, order.UserId)
				if tt.status != pb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
					assert.Equal(t, tt.status, order.Status)
				}
			}
		})
	}
}
