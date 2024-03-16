package service

import (
	"context"
	"database/sql"
	"github.com/gofrs/uuid/v5"
	"log"
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

	query := `
        WITH inserted AS (
            INSERT INTO orders (id, user_id, listing_id, status)
            VALUES ($1, $2, $3, $4)
            RETURNING *
        )
        SELECT inserted.id, inserted.listing_id, inserted.status, listings.name, listings.description, listings.uri
        FROM inserted
        INNER JOIN listings ON inserted.listing_id = listings.id
    `
	row := o.db.QueryRowContext(ctx, query, order.ID, order.UserID, order.ListingID, order.Status)

	var insertedOrder model.Order
	var name, description, tokenURI string
	if err := row.Scan(&insertedOrder.ID, &insertedOrder.ListingID, &insertedOrder.Status, &name, &description, &tokenURI); err != nil {
		return nil, err
	}

	log.Println("Order created: ", insertedOrder.ID.String())
	log.Println("Order listing: ", insertedOrder.ListingID.String())
	log.Println("Order status: ", insertedOrder.Status)
	log.Println("Order name: ", name)
	log.Println("Order description: ", description)
	log.Println("Order tokenURI: ", tokenURI)

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:          insertedOrder.ID.String(),
			ListingId:   insertedOrder.ListingID.String(),
			Name:        name,
			Description: description,
			TokenUri:    tokenURI,
			Status:      string(insertedOrder.Status),
		},
	}, nil
}

func (o orderService) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.OrdersResponse, error) {
	query := `
        SELECT o.id, o.user_id, o.listing_id, o.status, l.name, l.description, l.uri
        FROM orders o
        JOIN listings l ON o.listing_id = l.id
        WHERE o.user_id = $1
    `
	rows, err := o.db.QueryContext(ctx, query, req.GetUserId())
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order
	for rows.Next() {
		var order model.Order
		var name, description, tokenURI string
		err = rows.Scan(&order.ID, &order.UserID, &order.ListingID, &order.Status, &name, &description, &tokenURI)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &pb.Order{
			Id:          order.ID.String(),
			ListingId:   order.ListingID.String(),
			Name:        name,
			Description: description,
			TokenUri:    tokenURI,
			Status:      string(order.Status),
		})
	}

	return &pb.OrdersResponse{
		Orders: orders,
	}, nil
}

func (o orderService) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.OrderResponse, error) {
	var order model.Order
	var name, description, tokenURI string
	err := o.db.QueryRowContext(ctx, `
        SELECT o.id, o.user_id, o.listing_id, o.status, l.name, l.description, l.uri
        FROM orders o
        JOIN listings l ON o.listing_id = l.id
        WHERE o.id = $1
    `, req.GetId()).Scan(&order.ID, &order.UserID, &order.ListingID, &order.Status, &name, &description, &tokenURI)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:          order.ID.String(),
			ListingId:   order.ListingID.String(),
			Name:        name,
			Description: description,
			TokenUri:    tokenURI,
			Status:      string(order.Status),
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
