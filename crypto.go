package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

/* ============= Actions used for security in nanopm ============= */

/* Encrypt function */
func encrypt(records *[]byte, db_path *string, db_key *[]byte) {
	c, _ := aes.NewCipher(*db_key)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(*db_path, gcm.Seal(nonce, nonce, *records, nil), 0600)
}

/* Decrypt function. Returns json of the encrypted database */
func decrypt(key []byte, file string) []byte {
	ciphertext, _ := ioutil.ReadFile(file)
	c, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(c)
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("\nPassword is invalid")
		os.Exit(0)
	}
	return plaintext
}

/* Returns derived user password */
func getDerivedPassword(key *[]byte) *[]byte {
	salt := sha512.Sum512(*key)
	dk := pbkdf2.Key(*key, salt[:], 4096, 32, sha512.New)
	return &dk
}

/* Returns users password */
func createNewPassword() ([]byte, error) {
	db_pass := []byte(input("Enter the password for the database: "))
	if len(db_pass) < 5 {
		return nil, errors.New("newpass: password must be at least 5 symbols")
	}
	return db_pass, nil
}

/* Change password for the database */
func changeDatabasePassword(db_key *[]byte) {
	temp_key, pass_error := createNewPassword()
	if pass_error == nil {
		*db_key = *getDerivedPassword(&temp_key)
	} else {
		clearScreen()
		fmt.Println(pass_error)
		enterToContinue()
	}
}
