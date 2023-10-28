package betrens

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Topic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TopicName string             `json:"topicname" bson:"topicname"`
}
