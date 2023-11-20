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
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		user, err := module.GetUserFromID(objectId, db)
		if err != nil {
			t.Errorf("Error getting document: %v", err)
		} else {
			fmt.Println("Welcome bang:", user)
		}
	}

}

func TestAddTopic(t *testing.T) {
	var doc model.Topic
	doc.TopicName = "Erdito Nausha Adam2"
	doc.Source.Source = "yutube"
	doc.Source.Value = "https://twitter.com/erditonausha"
	// doc.Source.DateRange = "2021/01/01-2021/01/31"

	_id, err := module.AddTopic(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", _id)
	}
}

func TestGetTopic(t *testing.T) {
	var doc model.Topic
	id, err := primitive.ObjectIDFromHex("6549bae758638dd4f0137a11")
	doc.ID = id
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		doc, err := module.GetTopic(doc.ID, db)
		if err != nil {
			t.Errorf("Error getting document: %v", err)
		} else {
			fmt.Println("Welcome bang:", doc)
		}
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

func TestUpdateTopic(t *testing.T) {
	var doc model.Topic
	doc.TopicName = "dani"
	doc.Source.Source = "yutube"
	doc.Source.Value = "https://twitter.com/erditonausha"
	doc.Source.DateRange = "2021/01/01-2021/01/31"
	id, err := primitive.ObjectIDFromHex("653d4c011bdcf0c3ea14ee0a")
	doc.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah dengan id")
	} else {

		err = module.UpdateTopic(db, doc)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Println("Data berhasil diubah dengan id :", doc.ID)
		}
	}

}

func TestDeleteTopic(t *testing.T) {
	var doc model.Topic
	id, err := primitive.ObjectIDFromHex("653d4c011bdcf0c3ea14ee0a")
	doc.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil dihapus dengan id")
	} else {

		err = module.DeleteTopic(db, doc)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Println("Data berhasil dihapus dengan id :", doc.ID)
		}
	}
}

func TestEncrypt(t *testing.T) {
	var text = "daniaw"
	textEncrypt, err := module.EncryptString(text)
	if err != nil {
		t.Errorf("Error encrypting document: %v", err)
	} else {
		fmt.Println("Data berhasil dienkripsi :", textEncrypt)
	}
}

func TestDecrypt(t *testing.T) {
	var text = "daniaw"
	textEncrypt, err := module.EncryptString(text)
	if err != nil {
		t.Errorf("Error encrypting document: %v", err)
	} else {
		fmt.Println("Data berhasil dienkripsi :", textEncrypt)
	}
	textDecrypt, err := module.DecryptString(textEncrypt)
	if err != nil {
		t.Errorf("Error decrypting document: %v", err)
	} else {
		fmt.Println("Data berhasil didekripsi :", textDecrypt)
	}
}

func TestGenerateOTP(t *testing.T) {
	var email = ""
	otp, _ := module.OtpGenerate()
	var expiredAt = module.GenerateExpiredAt()
	var doc model.Otp
	doc.Email = email
	doc.OTP = otp
	doc.ExpiredAt = expiredAt
	fmt.Println(otp)
	fmt.Println(expiredAt)
}

func TestSendOTP(t *testing.T) {
	var email = ""
	otp, _ := module.OtpGenerate()
	var expiredAt = module.GenerateExpiredAt()
	var doc model.Otp
	doc.Email = email
	doc.OTP = otp
	doc.ExpiredAt = expiredAt
	fmt.Println(otp)
	fmt.Println(expiredAt)
	otp, err := module.SendOTP(db, "email@gmail.com")
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", otp)
	}

}
