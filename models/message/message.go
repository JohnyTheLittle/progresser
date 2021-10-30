package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dialogue struct {
	ID           string             `json:"_id" bson:"_id"`
	Participants [2]string          `json:"participants" bson:"participants"`
	Messages     []Message          `json:"messages" bson:"messages"`
	Header       string             `json:"header" bson:"header"`
	CreatedBy    primitive.ObjectID `json:"creator" bson:"creator"`
}

type Message struct {
	From string             `json:"from" bson:"from"`
	To   string             `json:"to" bson:"to"`
	Text string             `json:"message" bson:"message"`
	Date primitive.DateTime `json:"date" bson:"date"`
}
