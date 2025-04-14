package mappings

import (
	"time"

	"github.com/danikdaraboz/ecommerce/inventory-service/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductDB struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	Category    string             `bson:"category"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func ToDBModel(p *model.Product) (*ProductDB, error) {
	var dbID primitive.ObjectID
	var err error
	
	if p.ID != "" {
		dbID, err = primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return nil, err
		}
	}
	
	return &ProductDB{
		ID:          dbID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func ToDomainModel(p *ProductDB) *model.Product {
	return &model.Product{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}