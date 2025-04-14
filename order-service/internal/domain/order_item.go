package domain

type OrderItem struct {
    ProductID string `bson:"product_id"`
    Quantity  int32  `bson:"quantity"`
}