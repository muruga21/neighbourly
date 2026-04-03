package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName    string             `bson:"fullName" json:"fullName"`
	Email       string             `bson:"email" json:"email"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	Password    string             `bson:"password" json:"-"`
	Role        string             `bson:"role" json:"role"`
}
