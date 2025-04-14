package domain

type Order struct {
    ID     string       `bson:"_id"`
    UserID string       `bson:"user_id"`
    Items  []OrderItem  `bson:"items"`
    Status string       `bson:"status"`
    Total  float64      `bson:"total"`
}