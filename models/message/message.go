package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Sender    primitive.ObjectID `json:"sender" bson:"sender"`
	Addressee primitive.ObjectID `json:"addressee" bson:"addressee"`
	Text      string             `json:"message" bson:"message"`
	Date      primitive.DateTime `json:"date" bson:"date"`
}
