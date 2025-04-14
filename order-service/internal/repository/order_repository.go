package repository

import (
    "context"
    "order-service/internal/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
    collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
    return &OrderRepository{
        collection: db.Collection("orders"),
    }
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
    _, err := r.collection.InsertOne(ctx, order)
    return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
    var order domain.Order
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
    if err != nil {
        return nil, err
    }
    return &order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
    _, err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.M{"$set": bson.M{"status": status}},
    )
    return err
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string, page, pageSize int32) ([]*domain.Order, error) {
    skip := (page - 1) * pageSize
    cursor, err := r.collection.Find(
        ctx,
        bson.M{"user_id": userID},
        options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)),
    )
    if err != nil {
        return nil, err
    }
    var orders []*domain.Order
    if err := cursor.All(ctx, &orders); err != nil {
        return nil, err
    }
    return orders, nil
}