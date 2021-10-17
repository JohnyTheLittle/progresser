package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blacklist struct {
	Belongs       primitive.ObjectID   `json:"user" bson:"user"`
	ListOfPersons []primitive.ObjectID `json:"blocked_persons" bson:"blocked_persons"`
}
