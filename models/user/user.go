package models

type User struct {
	ID       string `bson:"_id" json:"_id"`
	Name     string `bson:"username" json:"name"`
	Password string `bson:"password" json:"password"`
	Email    string `bson:"email" json:"email"`
	URLName  string `bson:"url_name" json:"url_name"`
}
