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
	doc.Source.Source = "twitter"
	doc.Source.Value = "jokowi"
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
	id, err := primitive.ObjectIDFromHex("65657e1a29b3bad8b2bbb8e3")
	doc.ID = id
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		doc, _, err := module.GetTopic(doc.ID, db)
		if err != nil {
			t.Errorf("Error getting document: %v", err)
		} else {
			fmt.Println("Welcome bang:", doc)
			// fmt.Println("Welcome bang:", docs)
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

func TestCekOTP(t *testing.T) {
	var email = "email@gmail.com"
	otp := "6453"
	otp, err := module.VerifyOTP(db, email, otp)
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", otp)
	}
}

func TestUpdatePassword(t *testing.T) {
	var email = "email@gmail.com"
	otp := "6453"
	password := "daniaw"
	message, err := module.ResetPassword(db, email, otp, password)
	if err != nil {
		fmt.Println("Error sending otp: ", err)
	} else {
		fmt.Println("Data berhasil dikirim :", message)
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	phoneNumber := "62812345690"
	isValid, err := module.ValidatePhoneNumber(phoneNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isValid {
		fmt.Println("Phone number is valid.")
	} else {
		fmt.Println("Phone number is not valid.")
	}
}

func TestScrapSentimen(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("6566fa4d094a7eed963dd28f")
	topic := model.Topic{
		ID: id,
	}
	_, err := module.ScrapSentimen(db, topic)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// fmt.Println(docs)
}

func TestAddSetting(t *testing.T) {
	var doc model.Setting
	doc.ID = primitive.NewObjectID()
	doc.MaxTweeter = 100
	doc.MaxYoutube = 100
	doc.Twitter.UserName = "dani"
	doc.Twitter.Password = "dani"
	doc.Youtube.ApiKey = "dani"
	// doc.UserID = primitive.NewObjectID()

	_, err := module.AddSetting(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", doc.ID)
	}
}

func TestUpdateSetting(t *testing.T) {
	var doc model.Setting
	_id, err := primitive.ObjectIDFromHex("65644d6f66dc320b7b26f5e8")
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	}
	doc.ID = _id
	doc.MaxTweeter = 100
	doc.MaxYoutube = 100
	doc.Twitter.UserName = "dani"
	doc.Twitter.Password = "dani"
	doc.Youtube.ApiKey = "dani"
	err = module.UpdateSetting(db, doc)
	if err != nil {
		t.Errorf("Error updating document: %v", err)
	} else {
		fmt.Println("Data berhasil dirubah dengan id :", doc.ID)
	}
}

func TestGetSetting(t *testing.T) {
	var doc model.Setting
	_id, err := primitive.ObjectIDFromHex("656455b184d8d327072ba54b")
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	}
	doc.ID = _id
	setting, err := module.GetSetting(db, doc)

	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", setting)
	}
}

func TestGetAllSetting(t *testing.T) {
	var docs []model.Setting
	docs, err := module.GetAllSetting(db)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan id :", docs)
	}
	fmt.Println(docs)
}

func TestDeleteSetting(t *testing.T) {
	var doc model.Setting
	_id, err := primitive.ObjectIDFromHex("65644d6f66dc320b7b26f5e8")
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	}
	doc.ID = _id
	err = module.DeleteSetting(db, doc)
	if err != nil {
		t.Errorf("Error deleting document: %v", err)
	} else {
		fmt.Println("Data berhasil dihapus dengan id :", doc.ID)
	}
}
