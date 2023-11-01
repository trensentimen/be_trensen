package betrens

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	model "github.com/trensentimen/be_trensen/model"
	"github.com/whatsauth/watoken"
)

func GCFHandlerSignup(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUp(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + dataUser.Name
	return GCFReturnStruct(Response)
}

func GCFHandlerSignin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, status1, err := SignIn(conn, collectionname, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	// Response.Message = "Halo " + dataUser.Name
	tokenstring, err := watoken.Encode(dataUser.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email + " di Trensentimen" + strconv.FormatBool(status1)
		Response.Token = tokenstring
	}
	return GCFReturnStruct(Response)
}

func GCFHandlerGetAllTopic(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	var Response model.TopicResponse
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	topic, err := GetAllTopic(conn)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Selamat Datang " + dataUser.Email
	Response.Data = topic
	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
