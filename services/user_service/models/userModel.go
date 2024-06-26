package models

type UserModel struct {
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Phone string `json:"phone" bson:"phone"`
	Role string `json:"role" bson:"role"`
	District string `json:"district" bson:"district"`
	Address string `json:"address" bson:"address"`
}

func (u * UserModel) SetDefaults() {
	u.Name = "";
	u.Email = "";
	u.Password = "";
	u.Phone = "";
	u.Role = "";
	u.District = "";
	u.Address = "";
}