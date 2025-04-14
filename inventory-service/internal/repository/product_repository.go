package repository

import (
    "context"
    "inventory-service/internal/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
    collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
    return &ProductRepository{
        collection: db.Collection("products"),
    }
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
    _, err := r.collection.InsertOne(ctx, product)
    return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    var product domain.Product
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
    _, err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": product.ID},
        bson.M{"$set": product},
    )
    return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *ProductRepository) List(ctx context.Context, page, pageSize int32) ([]*domain.Product, error) {
    skip := (page - 1) * pageSize
    cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
    if err != nil {
        return nil, err
    }
    var products []*domain.Product
    if err := cursor.All(ctx, &products); err != nil {
        return nil, err
    }
    return products, nil
}