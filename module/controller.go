package betrens

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"

	"math/big"
	// "crypto/rand"
	// "encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/badoux/checkmail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/options"
	// "golang.org/x/crypto/argon2"

	model "github.com/trensentimen/be_trensen/model"
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://admin:admin@projectexp.pa7k8.gcp.mongodb.net"))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	// Define the regular expression pattern for numeric characters
	numericPattern := `^[0-9]+$`

	// Compile the numeric pattern
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	// Check if the phone number consists only of numeric characters
	if !numericRegexp.MatchString(phoneNumber) {
		return false, nil
	}

	// Define the regular expression pattern for "62" followed by 6 to 12 digits
	pattern := `^62\d{6,13}$`

	// Compile the regular expression
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	// Test if the phone number matches the pattern
	isValid := regexpPattern.MatchString(phoneNumber)

	return isValid, nil
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func GetOTPbyEmail(email string, db *mongo.Database) (doc model.Otp, err error) {
	collection := db.Collection("otp")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(db *mongo.Database, col string, id primitive.ObjectID, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("UpdatePresensi: %v\n", err)
		return
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

func SignUp(db *mongo.Database, col string, insertedDoc model.User) error {
	objectId := primitive.NewObjectID()

	if insertedDoc.Name == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.PhoneNumber == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	valid, _ := ValidatePhoneNumber(insertedDoc.PhoneNumber)
	if !valid {
		return fmt.Errorf("nomor telepon tidak valid")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}

	hash, _ := HashPassword(insertedDoc.Password)
	// insertedDoc.Password = hash
	user := bson.M{
		"_id":         objectId,
		"email":       insertedDoc.Email,
		"password":    hash,
		"role":        "user",
		"name":        insertedDoc.Name,
		"phonenumber": insertedDoc.PhoneNumber,
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return err
	}
	return nil
}

func SignIn(db *mongo.Database, col string, insertedDoc model.User) (user model.User, Status bool, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, false, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, false, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	if !CheckPasswordHash(insertedDoc.Password, existsDoc.Password) {
		return user, false, fmt.Errorf("password salah")
	}

	return existsDoc, true, nil
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func AddTopic(db *mongo.Database, doc model.Topic) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection("topic").InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func GetTopic(_id primitive.ObjectID, db *mongo.Database) (doc model.Topic, err error) {
	collection := db.Collection("topic")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetAllTopic(db *mongo.Database) (docs []model.Topic, err error) {
	collection := db.Collection("topic")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	return docs, nil
}

// update topic
func UpdateTopic(db *mongo.Database, doc model.Topic) (err error) {
	filter := bson.M{"_id": doc.ID}
	result, err := db.Collection("topic").UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Printf("UpdateTopic: %v\n", err)
		return
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

// delete topic
func DeleteTopic(db *mongo.Database, doc model.Topic) error {
	collection := db.Collection("topic")
	filter := bson.M{"_id": doc.ID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", doc.ID, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", doc.ID)
	}

	return nil
}

func EncryptString(text string) (string, error) {
	key := []byte("a16bytekeyforAES")
	// text := "Hello, Golang!"

	encryptedText, err := encrypt(text, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return "", err
	}

	return encryptedText, nil
}

func DecryptString(encryptedText string) (string, error) {
	key := []byte("a16bytekeyforAES")

	decryptedText, err := decrypt(encryptedText, key)
	if err != nil {
		// fmt.Println("Error decrypting:", err)
		return "", err
	}

	// fmt.Println("Decrypted:", decryptedText)
	return decryptedText, nil
}

func OtpGenerate() (string, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	// Format the random number as a 4-digit string
	otp := fmt.Sprintf("%04d", randomNumber)

	return otp, nil
}

func GenerateExpiredAt() int64 {
	currentTime := time.Now()

	// Add 5 minutes
	newTime := currentTime.Add(5 * time.Minute)
	return newTime.Unix()
}

func SendOTP(db *mongo.Database, email string) (string, error) {
	// GET OTP
	otp, _ := OtpGenerate()

	// GET EXPIRED AT
	expiredAt := GenerateExpiredAt()

	// get user by email
	existsDoc, err := GetUserFromEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("email tidak ditemukan1")
	}
	if existsDoc.Email == "" {
		return "", fmt.Errorf("email tidak ditemukan2")
	}

	// save otp to db
	// objectId := primitive.NewObjectID()
	otpDoc := bson.M{
		// "_id":       objectId,
		"email":     email,
		"otp":       otp,
		"expiredat": expiredAt,
		"status":    false,
	}

	// get otp by email
	_, err = GetOTPbyEmail(email, db)

	if err != nil {
		if err.Error() == "email tidak ditemukan" {
			// return "", fmt.Errorf("error getting OTP from email: %s", err.Error())
			// insert new OTP
			_, err = db.Collection("otp").InsertOne(context.Background(), otpDoc)
			if err != nil {
				return "", fmt.Errorf("error inserting OTP: %s", err.Error())
			}
			return otp, nil
		} else {
			return "", fmt.Errorf("error Get OTP: %s", err.Error())
		}
	} else {
		// update existing OTP
		filter := bson.M{"email": email}
		update := bson.M{"$set": otpDoc}
		_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)
		if err != nil {
			return "", fmt.Errorf("error updating OTP: %s", err.Error())
		}
	}

	// postapi
	url := "https://api.wa.my.id/api/send/message/text"

	// Data yang akan dikirimkan dalam format JSON
	jsonStr := []byte(`{
        "to": "` + existsDoc.PhoneNumber + `",
        "isgroup": false,
        "messages": "kode Otp akun trensentimen.my.id atas nama *` + email + `* adalah *` + otp + `*"
    }`)

	// Membuat permintaan HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Menambahkan header ke permintaan
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", "v4.public.eyJleHAiOiIyMDIzLTEyLTE2VDE3OjQ3OjQ3KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0xNlQxNzo0Nzo0NyswNzowMCIsImlkIjoiNjI4NTcwMzMwNTE2MyIsIm5iZiI6IjIwMjMtMTEtMTZUMTc6NDc6NDcrMDc6MDAifXlYzCjMwUnUHhdyWpcQyq33tOKlhJIWHzBr5Zq2PgmYxjeghbWqkS1QUH7ojfzPYd1fIaWOHnoE29zbE-v_tQk")
	req.Header.Set("Content-Type", "application/json")

	// Melakukan permintaan HTTP POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Menampilkan respons dari server
	fmt.Println("Response Status:", resp.Status)
	return "success", nil
}

func VerifyOTP(db *mongo.Database, email, otp string) (string, error) {
	// get otp by email
	otpDoc, err := GetOTPbyEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("error Get OTP: %s", err.Error())
	}

	// check otp
	if otpDoc.OTP != otp {
		return "", fmt.Errorf("otp tidak valid")
	}

	// check expired at
	if otpDoc.ExpiredAt < time.Now().Unix() {
		return "", fmt.Errorf("otp telah kadaluarsa")
	}

	//update otp
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"status": true}}
	_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)

	if err != nil {
		return "", fmt.Errorf("error updating OTP: %s", err.Error())
	}

	return otp, nil
}

func ResetPassword(db *mongo.Database, email, otp, password string) (string, error) {
	// get user by email
	existsDoc, err := GetUserFromEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("email tidak ditemukan1")
	}
	if existsDoc.Email == "" {
		return "", fmt.Errorf("email tidak ditemukan2")
	}

	// check otp
	docOtp, err := GetOTPbyEmail(email, db)
	if err != nil {
		return "", fmt.Errorf("error Get OTP: %s", err.Error())
	}
	if docOtp.OTP != otp || !docOtp.Status {
		return "", fmt.Errorf("otp tidak valid")
	}

	// hash password
	hash, _ := HashPassword(password)

	// update password
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": hash}}
	_, err = db.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", fmt.Errorf("error updating password: %s", err.Error())
	}

	// update otp
	filter = bson.M{"email": email}
	update = bson.M{"$set": bson.M{"status": false}}
	_, err = db.Collection("otp").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", fmt.Errorf("error updating password: %s", err.Error())
	}

	return "success", nil
}

func InsertManyDocs(db *mongo.Database, col string, dataTopics []model.DataTopics) (insertedIDs []primitive.ObjectID, err error) {
	var interfaces []interface{}
	for _, topic := range dataTopics {
		interfaces = append(interfaces, topic)
	}
	result, err := db.Collection(col).InsertMany(context.Background(), interfaces)
	if err != nil {
		return insertedIDs, fmt.Errorf("kesalahan server: insert")
	}
	for _, id := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID))
	}
	return insertedIDs, nil
}

func ScrapSentimen(db *mongo.Database, topic model.Topic) (docs []model.DataTopics, err error) {

	if topic.Source.Source == "youtube" {
		return docs, fmt.Errorf("fitur sedang dikembangkan")
		// docs, err = CrawlingYoutube(topic)
		// if err != nil {
		// 	return docs, fmt.Errorf("error CrawlingYoutube: %s", err.Error())
		// }
	} else if topic.Source.Source == "twitter" {
		docs, err = CrawlingTweet(topic)
		if err != nil {
			return docs, fmt.Errorf("error CrawlingTweet: %s", err.Error())
		}
	} else {
		return docs, fmt.Errorf("source tidak ditemukan")
	}

	// dataTopics, err := CrawlingTweet(topic)
	// if err != nil {
	// 	return docs, fmt.Errorf("error CrawlingTweet: %s", err.Error())
	// }

	// insert data to db
	_, err = InsertManyDocs(db, "datatopics", docs)
	if err != nil {
		return docs, fmt.Errorf("error insert data: %s", err.Error())
	}

	// for _, data := range dataTopics {
	// 	_, err = InsertOneDoc(db, "datatopics", data)
	// 	if err != nil {
	// 		return docs, fmt.Errorf("error insert data: %s", err.Error())
	// 	}
	// }

	return docs, nil
}
