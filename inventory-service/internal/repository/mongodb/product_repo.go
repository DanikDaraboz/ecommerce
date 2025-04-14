package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/danikdaraboz/ecommerce/inventory-service/internal/domain"
	"github.com/danikdaraboz/ecommerce/inventory-service/internal/domain/model"
	"github.com/danikdaraboz/ecommerce/inventory-service/internal/repository/mongodb/mappings"
)

// Ensure ProductRepository implements domain.ProductRepository
var _ domain.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewProductRepository(uri, dbName string) (domain.ProductRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &ProductRepository{
		client:     client,
		collection: client.Database(dbName).Collection("products"),
	}, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	if product == nil {
		return nil, errors.New("product cannot be nil")
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

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
	if id == "" {
		return nil, domain.ErrInvalidID
	}

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
	if product == nil {
		return errors.New("product cannot be nil")
	}

	product.UpdatedAt = time.Now()

	dbProduct, err := mappings.ToDBModel(product)
	if err != nil {
		return err
	}

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": dbProduct.ID},
		dbProduct,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context, skip, limit int64) ([]*model.Product, error) {
	if skip < 0 || limit <= 0 {
		return nil, errors.New("invalid pagination parameters")
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1}) // Sort by newest first

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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}