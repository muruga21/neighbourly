package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Service struct {
	ID    string `bson:"id" json:"id"`
	Icon  string `bson:"icon" json:"icon"`
	Title string `bson:"title" json:"title"`
	Desc  string `bson:"desc" json:"desc"`
	Price string `bson:"price" json:"price"`
}

type Review struct {
	ID     string `bson:"id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
	Rating int    `bson:"rating" json:"rating"`
	Text   string `bson:"text" json:"text"`
}

type Location struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName          string             `bson:"fullName" json:"fullName"`
	Email             string             `bson:"email" json:"email"`
	PhoneNumber       string             `bson:"phoneNumber" json:"phoneNumber"`
	Password          string             `bson:"password" json:"-"`
	Role              string             `bson:"role" json:"role"`
	IsProfileComplete bool               `bson:"isProfileComplete" json:"isProfileComplete"`
	Avatar            string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Title             string             `bson:"title,omitempty" json:"title,omitempty"`
	Experience        string             `bson:"experience,omitempty" json:"experience,omitempty"`
	BasePrice         string             `bson:"basePrice,omitempty" json:"basePrice,omitempty"`
	Status            string             `bson:"status,omitempty" json:"status,omitempty"`
	Location          *Location          `bson:"location,omitempty" json:"location,omitempty"`
	Services          []Service          `bson:"services,omitempty" json:"services,omitempty"`
	Reviews           []Review           `bson:"reviews,omitempty" json:"reviews,omitempty"`
}
