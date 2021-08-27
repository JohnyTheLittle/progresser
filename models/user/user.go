package user

type User struct {
	Name     string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}
