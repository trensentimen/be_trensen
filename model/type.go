package betrens

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Password    string             `json:"password" bson:"password"`
	Email       string             `json:"email" bson:"email"`
	Role        string             `json:"role,omitempty" bson:"role,omitempty"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Topic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TopicName string             `json:"topicname" bson:"topicname"`
	Source    Sources            `json:"source" bson:"source"`
	Status    string             `json:"status" bson:"status"` // Drafting, Inputting, Analyzing
}

type DataTopics struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TopicId  primitive.ObjectID `json:"topicid,omitempty" bson:"topicid,omitempty"`
	Text     string             `json:"text" bson:"text"`
	Sentimen string             `json:"sentimen" bson:"sentimen"` // positive/negative/neutral
	Source   string             `json:"source" bson:"source"`     // youtube/twitter
	Date     int64              `json:"date" bson:"date"`
}

type TopicResponse struct {
	Status  bool    `json:"status" bson:"status"`
	Message string  `json:"message,omitempty" bson:"message,omitempty"`
	Data    []Topic `json:"data" bson:"data"`
}

type Sources struct {
	Source    string `json:"source" bson:"source"` // youtube/twitter
	Value     string `json:"value" bson:"value"`   // link
	DateRange string `json:"date_range" bson:"date_range"`
}

type Otp struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	OTP       string             `json:"otp" bson:"otp"`
	ExpiredAt int64              `json:"expiredat" bson:"expiredat"`
	Status    bool               `json:"status" bson:"status"`
}

type Response struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type ResetPassword struct {
	Email    string `json:"email" bson:"email"`
	OTP      string `json:"otp" bson:"otp"`
	Password string `json:"password" bson:"password"`
}

// type Setting struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
// 	SettingName string             `json:"settingname" bson:"settingname"`
// 	SettingType string             `json:"settingtype" bson:"settingtype"`
// 	Setting     string             `json:"setting" bson:"setting"`
// }

type Setting struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	// UserID     primitive.ObjectID `bson:"userid,omitempty" json:"userid,omitempty"`
	MaxTweeter int     `json:"maxtweet" bson:"maxtweet"`
	MaxYoutube int     `json:"maxyoutube" bson:"maxyoutube"`
	Twitter    Twitter `json:"twitter" bson:"twitter"`
	Youtube    Youtube `json:"youtube" bson:"youtube"`
}

type Twitter struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Youtube struct {
	ApiKey string `json:"apikey" bson:"apikey"`
}
