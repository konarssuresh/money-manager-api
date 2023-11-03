package helper

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Encrypt(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return string(hashedPassword), err
	}

	return string(hashedPassword), nil

}

func ValidateEncryption(encryptedStr string, inputString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedStr), []byte(inputString))

	if err != nil {
		return false
	} else {
		return true
	}
}

func GetTodaysDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func GetCurrentTimeStamp() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05.999999-07")
}
