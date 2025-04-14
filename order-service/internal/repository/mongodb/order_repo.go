package mongodb

import (
	"context"
	"time"

	"github.com/danikdaraboz/ecommerce/order-service/internal/domain"
	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
	"github.com/danikdaraboz/ecommerce/order-service/internal/repository/mongodb/mappings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepositoryMongo struct {
	collection *mongo.Collection
}

func NewOrderRepositoryMongo(collection *mongo.Collection) domain.OrderRepository {
	return &OrderRepositoryMongo{collection: collection}
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	dbModel, err := mappings.ToDBModel(order)
	if err != nil {
		return nil, err
	}

	res, err := r.collection.InsertOne(ctx, dbModel)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	order.ID = id
	return order, nil
}

func (r *OrderRepositoryMongo) FindByID(ctx context.Context, id string) (*model.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidOrder
	}

	var dbModel mappings.OrderDB
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&dbModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return mappings.ToDomainModel(&dbModel), nil
}

func (r *OrderRepositoryMongo) FindByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*model.Order

	for cursor.Next(ctx) {
		var dbModel mappings.OrderDB
		if err := cursor.Decode(&dbModel); err != nil {
			return nil, err
		}
		order := mappings.ToDomainModel(&dbModel)
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryMongo) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidOrder
	}

	update := bson.M{
		"$set": bson.M{
			"status":     string(status),
			"updated_at": time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

func (r *OrderRepositoryMongo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidOrder
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}
