package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, user *User) (string, error) {
	// get user form email
	dbUser := new(User)
	result := db.Where("email = ?", user.Email).First(dbUser)

	if result.Error != nil {
		return "", result.Error
	}
	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

	if err != nil {
		return "", err
	}
	// if pass create jwt token
	jwtSecretKey := "TestSecret"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = dbUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}
