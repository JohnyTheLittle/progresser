package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dialogue struct {
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
	Messages     []Message            `json:"messages" bson:"messages"`
	Header       string               `json:"header" bson:"header"`
	CreatedBy    primitive.ObjectID   `json:"creator" bson:"creator"`
}

type Message struct {
	Sender    primitive.ObjectID `json:"sender" bson:"sender"`
	Addressee primitive.ObjectID `json:"addressee" bson:"addressee"`
	Text      string             `json:"message" bson:"message"`
	Date      primitive.DateTime `json:"date" bson:"date"`
}
