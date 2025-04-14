package domain

type Product struct {
    ID          string  `bson:"_id"`
    Name        string  `bson:"name"`
    Description string  `bson:"description"`
    Price       float64 `bson:"price"`
    Stock       int32   `bson:"stock"`
    Category    string  `bson:"category"`
}