package core

type User struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Age      string `json:"age" bson:"age"`
	Password string `json:"-" bson:"password"`
}
