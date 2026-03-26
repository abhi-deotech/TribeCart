package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tribecart/proto/tribecart/v1"
)

var (
	orderClient pb.OrderServiceClient
)

// Request/Response types
type (
	registerRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	createOrderRequest struct {
		UserID           string          `json:"user_id"`
		Items            []*pb.OrderItem `json:"items"`
		ShippingAddress  string          `json:"shipping_address_id"`
		PaymentMethod    string          `json:"payment_method"`
		PaymentID        string          `json:"payment_id"`
	}

	updateOrderStatusRequest struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number,omitempty"`
		Carrier        string `json:"carrier,omitempty"`
	}

	cancelOrderRequest struct {
		Reason        string `json:"reason"`
		RefundPayment bool   `json:"refund_payment"`
	}

	listOrdersRequest struct {
		PageSize  int32     `json:"page_size"`
		PageToken string    `json:"page_token"`
		UserID    string    `json:"user_id"`
		Status    string    `json:"status"`
		FromDate  time.Time `json:"from_date"`
		ToDate    time.Time `json:"to_date"`
	}

	errorResponse struct {
		Error string `json:"error"`
	}
)

// Handlers
func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func registerUser(client pb.UserServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleError(w, err, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &pb.CreateUserRequest{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  req.Password,
		}

		res, err := client.CreateUser(ctx, grpcReq)
		if err != nil {
			handleGRPCError(w, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, res)
	}
}

// Order Handlers
func createOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert payment method string to enum
	paymentMethod, ok := pb.PaymentMethod_value[req.PaymentMethod]
	if !ok {
		handleError(w, nil, "Invalid payment method", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CreateOrderRequest{
		UserId:            req.UserID,
		Items:             req.Items,
		ShippingAddressId: req.ShippingAddress,
		PaymentMethod:     pb.PaymentMethod(paymentMethod),
		PaymentId:         req.PaymentID,
	}

	order, err := orderClient.CreateOrder(ctx, grpcReq)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	if orderID == "" {
		handleError(w, nil, "Order ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req := &pb.GetOrderRequest{
		Id:            orderID,
		IncludeItems:  true,
		IncludeAddress: true,
	}

	order, err := orderClient.GetOrder(ctx, req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

func listOrders(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	pageSize, _ := strconv.ParseInt(query.Get("page_size"), 10, 32)
	status, _ := pb.OrderStatus_value[query.Get("status")]

	req := &pb.ListOrdersRequest{
		PageSize: int32(pageSize),
		PageToken: query.Get("page_token"),
		UserId:    query.Get("user_id"),
		Status:    pb.OrderStatus(status),
	}

	// Parse date filters
	if fromDate := query.Get("from_date"); fromDate != "" {
		if t, err := time.Parse(time.RFC3339, fromDate); err == nil {
			req.FromDate = timestamppb.New(t)
		}
	}
	if toDate := query.Get("to_date"); toDate != "" {
		if t, err := time.Parse(time.RFC3339, toDate); err == nil {
			req.ToDate = timestamppb.New(t)
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := orderClient.ListOrders(ctx, req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	if orderID == "" {
		handleError(w, nil, "Order ID is required", http.StatusBadRequest)
		return
	}

	var req updateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	status, ok := pb.OrderStatus_value[req.Status]
	if !ok {
		handleError(w, nil, "Invalid order status", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.UpdateOrderStatusRequest{
		Id:            orderID,
		Status:        pb.OrderStatus(status),
		TrackingNumber: req.TrackingNumber,
		Carrier:       req.Carrier,
	}

	order, err := orderClient.UpdateOrderStatus(ctx, grpcReq)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	if orderID == "" {
		handleError(w, nil, "Order ID is required", http.StatusBadRequest)
		return
	}

	var req cancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := orderClient.CancelOrder(ctx, &pb.CancelOrderRequest{
		Id:           orderID,
		Reason:       req.Reason,
		RefundPayment: req.RefundPayment,
	})

	if err != nil {
		handleGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUserOrderHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		handleError(w, nil, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	pageSize, _ := strconv.ParseInt(query.Get("page_size"), 10, 32)
	status, _ := pb.OrderStatus_value[query.Get("status")]

	req := &pb.GetOrderHistoryRequest{
		UserId:    userID,
		PageSize:  int32(pageSize),
		PageToken: query.Get("page_token"),
		Status:    pb.OrderStatus(status),
	}

	// Parse date filters
	if fromDate := query.Get("from_date"); fromDate != "" {
		if t, err := time.Parse(time.RFC3339, fromDate); err == nil {
			req.FromDate = timestamppb.New(t)
		}
	}
	if toDate := query.Get("to_date"); toDate != "" {
		if t, err := time.Parse(time.RFC3339, toDate); err == nil {
			req.ToDate = timestamppb.New(t)
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := orderClient.GetOrderHistory(ctx, req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		// Use protojson for protobuf messages to handle them correctly
		if msg, ok := data.(proto.Message); ok {
			marshaler := protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			}
			jsonBytes, err := marshaler.Marshal(msg)
			if err != nil {
				handleError(w, err, "Error encoding response", http.StatusInternalServerError)
				return
			}
			w.Write(jsonBytes)
			return
		}

		// Regular JSON encoding for non-protobuf messages
		if err := json.NewEncoder(w).Encode(data); err != nil {
			handleError(w, err, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func handleError(w http.ResponseWriter, err error, message string, statusCode int) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	} else {
		log.Println(message)
	}

	respondWithJSON(w, statusCode, errorResponse{
		Error: message,
	})
}

func handleGRPCError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			handleError(w, err, st.Message(), http.StatusNotFound)
		case codes.InvalidArgument:
			handleError(w, err, st.Message(), http.StatusBadRequest)
		case codes.AlreadyExists:
			handleError(w, err, st.Message(), http.StatusConflict)
		case codes.PermissionDenied:
			handleError(w, err, st.Message(), http.StatusForbidden)
		case codes.Unauthenticated:
			handleError(w, err, "Unauthenticated", http.StatusUnauthorized)
		default:
			handleError(w, err, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	handleError(w, err, "Internal server error", http.StatusInternalServerError)
}

func main() {
	// gRPC client for User Service
	userAddr := os.Getenv("USER_SERVICE_ADDR")
	log.Printf("Connecting to User Service at: %s", userAddr)
	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := pb.NewUserServiceClient(userConn)

	// gRPC client for Order Service
	orderAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if orderAddr == "" {
		orderAddr = "orders-service:8080"
	}
	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to order service: %v", err)
	}
	defer orderConn.Close()
	orderClient = pb.NewOrderServiceClient(orderConn)

	// Create router
	r := mux.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Apply CORS to all routes
	r.Use(corsMiddleware.Handler)

	// Health check endpoint
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// User service routes
	userRouter := r.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("/register", registerUser(userClient)).Methods("POST")

	// Order service routes
	orderRouter := r.PathPrefix("/api/v1/orders").Subrouter()
	orderRouter.HandleFunc("", createOrder).Methods("POST")
	orderRouter.HandleFunc("", listOrders).Methods("GET")
	orderRouter.HandleFunc("/{id}", getOrder).Methods("GET")
	orderRouter.HandleFunc("/{id}/cancel", cancelOrder).Methods("POST")
	orderRouter.HandleFunc("/{id}/status", updateOrderStatus).Methods("PATCH")
	orderRouter.HandleFunc("/user/{user_id}", getUserOrderHistory).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API Gateway server started on %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
