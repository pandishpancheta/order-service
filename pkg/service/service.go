package service

import (
	"context"
	"database/sql"
	"github.com/gofrs/uuid/v5"
	"order-service/pkg/config"
	"order-service/pkg/model"
	"order-service/pkg/pb"
)

type Services interface {
	OrderService
}

type services struct {
	OrderService
}

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.NewOrderRequest) (*pb.OrderResponse, error)
	GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.OrdersResponse, error)
	GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.OrderResponse, error)
	UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.EmptyResponse, error)
}

type orderService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewOrderService(db *sql.DB, cfg *config.Config) OrderService {
	return services{
		OrderService: &orderService{
			db:  db,
			cfg: cfg,
		},
	}
}

func (o orderService) CreateOrder(ctx context.Context, req *pb.NewOrderRequest) (*pb.OrderResponse, error) {
	var order model.Order
	var err error

	order.ID, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}

	order.UserID, err = uuid.FromString(req.GetUserId())
	if err != nil {
		return nil, err
	}

	order.ListingID, err = uuid.FromString(req.GetListingId())
	if err != nil {
		return nil, err
	}

	order.Status = model.Pending
	order.TokenURI = req.GetTokenUri()

	_, err = o.db.ExecContext(ctx, "INSERT INTO orders (id, user_id, listing_id, status, token_uri) VALUES ($1, $2, $3, $4, $5)", order.ID, order.UserID, order.ListingID, order.Status, order.TokenURI)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:        order.ID.String(),
			UserId:    order.UserID.String(),
			ListingId: order.ListingID.String(),
			Status:    string(order.Status),
			TokenUri:  order.TokenURI,
		},
	}, nil
}

func (o orderService) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.OrdersResponse, error) {
	rows, err := o.db.QueryContext(ctx, "SELECT * FROM orders WHERE user_id = $1", req.GetUserId())
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.ID, &order.UserID, &order.ListingID, &order.Status, &order.TokenURI)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &pb.Order{
			Id:        order.ID.String(),
			UserId:    order.UserID.String(),
			ListingId: order.ListingID.String(),
			Status:    string(order.Status),
			TokenUri:  order.TokenURI,
		})
	}

	return &pb.OrdersResponse{
		Orders: orders,
	}, nil
}

func (o orderService) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.OrderResponse, error) {
	var order model.Order
	err := o.db.QueryRowContext(ctx, "SELECT * FROM orders WHERE id = $1", req.GetId()).Scan(&order.ID, &order.UserID, &order.ListingID, &order.Status, &order.TokenURI)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:        order.ID.String(),
			UserId:    order.UserID.String(),
			ListingId: order.ListingID.String(),
			Status:    string(order.Status),
			TokenUri:  order.TokenURI,
		},
	}, nil
}

func (o orderService) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.EmptyResponse, error) {
	_, err := o.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", req.GetStatus(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
