package betrens

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// func TestUpdateGetData(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "trensentimen")
// 	datagedung := GetAllBangunanLineString(mconn, "trensentimen")
// 	fmt.Println(datagedung)
// }

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("private key")
	fmt.Println(privateKey)
	fmt.Println("public key")
	fmt.Println(publicKey)
	fmt.Println("hasil")
	hasil, err := watoken.Encode("Trensentimen", privateKey)
	fmt.Println(hasil, err)
}

func TestValidateToken(t *testing.T) {
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEwLTE4VDE0OjIwOjE4WiIsImlhdCI6IjIwMjMtMTAtMThUMTI6MjA6MThaIiwiaWQiOiJkYW5pIiwibmJmIjoiMjAyMy0xMC0xOFQxMjoyMDoxOFoiffzDySc-bv3-c0Nv8xbJmWfjEntNupxaEFfom7WMocHgd8kRQyLg5DPs1poyv9RCJouac3CPfMI2eG1TS--zmQA" // Gantilah dengan token PASETO yang sesuai
	publicKey := "f48bd58cb3b3972d05bb9303b15ce9b83f4fcb9c871d1b05906f2fec20620ea0"
	payload, _err := watoken.Decode(publicKey, tokenstring)
	if _err != nil {
		fmt.Println("expire token", _err)
	} else {
		fmt.Println(payload.Id)
		fmt.Println(payload.Nbf)
		fmt.Println(payload.Iat)
		fmt.Println(payload.Exp)
	}

}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "trensentimen")

	var userdata User
	userdata.Username = "dani"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "trensentimen")
	var userdata User
	userdata.Username = "dani"
	userdata.Password = "secret"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "trensentimen")
	var userdata User
	userdata.Username = "dani2"
	userdata.Role = "Admin"
	userdata.Password = "secret"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestGCFPostHandler(t *testing.T) {

	// Membuat body request sebagai string
	requestBody := `{"username": "dani", "password": "secret"}`

	// Membuat objek http.Request
	r := httptest.NewRequest("POST", "https://contoh.com/path", strings.NewReader(requestBody))
	r.Header.Set("Content-Type", "application/json")

	resp := GCFPostHandler("PASETOPRIVATEKEY", "MONGOSTRING", "trensentimen", "user", r)
	fmt.Println(resp)
}
