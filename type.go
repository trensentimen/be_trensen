package betrens

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
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
	Range     string             `json:"range" bson:"range"`
	Source    Sources            `json:"source" bson:"source"`
}

type Sources struct {
	Name  string `json:"source" bson:"source"` // youtube/twitter
	Value string `json:"value" bson:"value"`   // link
}
