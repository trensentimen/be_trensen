package betrens

import (
	"fmt"
	"testing"

	model "github.com/trensentimen/be_trensen/model"
	module "github.com/trensentimen/be_trensen/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "trensentimen")

func TestSignUp(t *testing.T) {
	var doc model.User
	doc.Name = "Erdito Nausha Adam"
	doc.Email = "erdito2@gmail.com"
	doc.Password = "fghjkliow"
	doc.Role = "admin"

	err := module.SignUp(db, "user", doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", doc.Name)
	}
}

func TestLogIn(t *testing.T) {
	var doc model.User
	doc.Email = "erdito2@gmail.com"
	doc.Password = "fghjkliow"
	user, Status, err := module.SignIn(db, "user", doc)
	fmt.Println("Status :", Status)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome bang:", user)
	}
}

func TestGetAllDocs(t *testing.T) {
	var docs []model.User
	docs = module.GetAllDocs(db, "user", docs).([]model.User)
	fmt.Println(docs)
}

func TestGetUserFromID(t *testing.T) {
	// var doc model.User
	// doc.ID = "653d3367e56f0084ac013212"
	id := "653d3367e56f0084ac013212"
	objectId, err := primitive.ObjectIDFromHex(id)
	user, err := module.GetUserFromID(objectId, db)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome bang:", user)
	}
}

func TestAddTopic(t *testing.T) {
	var doc model.Topic
	doc.TopicName = "Erdito Nausha Adam2"
	doc.Source.Name = "yutube"
	doc.Source.Value = "https://twitter.com/erditonausha"
	// doc.Source.DateRange = "2021/01/01-2021/01/31"

	_id, err := module.AddTopic(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", _id)
	}
}

func TestGetAllTopic(t *testing.T) {
	var docs []model.Topic
	docs, err := module.GetAllTopic(db)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", docs)
	}
	fmt.Println(docs)
}
