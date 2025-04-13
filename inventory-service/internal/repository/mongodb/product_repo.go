package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"inventoryservice/internal/domain"
	"inventoryservice/internal/domain/model"
	"inventoryservice/internal/repository/mongodb/mappings"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(uri, dbName string) (*ProductRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &ProductRepository{
		collection: client.Database(dbName).Collection("products"),
	}, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	dbProduct, err := mappings.ToDBModel(product)
	if err != nil {
		return nil, err
	}

	res, err := r.collection.InsertOne(ctx, dbProduct)
	if err != nil {
		return nil, err
	}

	dbProduct.ID = res.InsertedID.(primitive.ObjectID)
	return mappings.ToDomainModel(dbProduct), nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var dbProduct mappings.ProductDB
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&dbProduct)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return mappings.ToDomainModel(&dbProduct), nil
}

func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	dbProduct, err := mappings.ToDBModel(product)
	if err != nil {
		return err
	}

	_, err = r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": dbProduct.ID},
		dbProduct,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *ProductRepository) List(ctx context.Context, skip, limit int64) ([]*model.Product, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	for cursor.Next(ctx) {
		var dbProduct mappings.ProductDB
		if err := cursor.Decode(&dbProduct); err != nil {
			return nil, err
		}
		products = append(products, mappings.ToDomainModel(&dbProduct))
	}

	return products, nil
}

func (r *ProductRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.collection.Database().Client().Disconnect(ctx)
}