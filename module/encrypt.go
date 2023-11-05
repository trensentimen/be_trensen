package betrens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Anda dapat melakukan enkripsi dan dekripsi teks dalam bahasa pemrograman Go (Golang) menggunakan berbagai metode kriptografi yang tersedia. Salah satu cara yang umum digunakan adalah dengan menggunakan AES (Advanced Encryption Standard) untuk enkripsi dan dekripsi teks. Berikut adalah contoh sederhana cara melakukan enkripsi dan dekripsi teks menggunakan AES di Golang

// Pastikan untuk mengganti nilai key sesuai dengan kunci enkripsi yang ingin Anda gunakan. Dalam contoh di atas, kami menggunakan AES dengan mode CFB (Cipher Feedback) untuk enkripsi dan dekripsi. Anda dapat menyesuaikan metode dan mode enkripsi sesuai kebutuhan Anda.

func encrypt(text string, key []byte) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
