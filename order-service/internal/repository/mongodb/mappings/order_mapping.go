package mappings

import (
	"time"

	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItemDB struct {
	ProductID string  `bson:"product_id"`
	Quantity  int     `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type OrderDB struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	Items     []OrderItemDB      `bson:"items"`
	Total     float64            `bson:"total"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func ToDBModel(order *model.Order) (*OrderDB, error) {
	if order == nil {
		return nil, nil
	}

	var dbID primitive.ObjectID
	if order.ID != "" {
		var err error
		dbID, err = primitive.ObjectIDFromHex(order.ID)
		if err != nil {
			return nil, err
		}
	}

	var items []OrderItemDB
	for _, item := range order.Items {
		items = append(items, OrderItemDB{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &OrderDB{
		ID:        dbID,
		UserID:    order.UserID,
		Items:     items,
		Total:     order.Total,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

func ToDomainModel(order *OrderDB) *model.Order {
	if order == nil {
		return nil
	}

	var items []model.OrderItem
	for _, item := range order.Items {
		items = append(items, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &model.Order{
		ID:        order.ID.Hex(),
		UserID:    order.UserID,
		Items:     items,
		Total:     order.Total,
		Status:    model.OrderStatus(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}