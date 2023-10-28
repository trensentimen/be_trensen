package betrens

import (
	"fmt"
	"testing"

	model "github.com/trensentimen/be_trensen/model"
	module "github.com/trensentimen/be_trensen/module"
)

var db = module.MongoConnect("MONGOSTRING", "trensentimen")

func TestSignUp(t *testing.T) {
	var doc model.User
	doc.Name = "Erdito Nausha Adam"
	doc.Email = "erdito2@gmail.com"
	doc.Password = "fghjkliow"

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
	user, err := module.SignIn(db, "user", doc)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome bang:", user)
	}
}
