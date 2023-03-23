package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(100);unique_index" json:"email"`
	FirstName string `gorm:"size:100;not null"              json:"firstname"`
	LastName  string `gorm:"size:100;not null"              json:"lastname"`
	Password  string `gorm:"size:100;not null"              json:"password"`
	Dummy     bool   `gorm:"default:false"`
}

func (u *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}

	err := db.Debug().Table("users").Where("email = ?", u.Email).First(account).Error

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	err := db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) ValidateUser(action string) error {

	switch action {
	case "login":
		{
			if u.Email == "" {
				return errors.New("Email is required")
			}
			if u.Password == "" {
				return errors.New("Password is required")
			}
			return nil
		}
	default:
		{
			if u.FirstName == "" {
				return errors.New("FirstName is required")
			}
			if u.LastName == "" {
				return errors.New("LastName is required")
			}
			if u.Email == "" {
				return errors.New("Email is required")
			}
			if u.Password == "" {
				return errors.New("Password is required")
			}
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid Email")
			}
			return nil
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password.
	return string(bytes), err
}

func (u *User) BeforeSave(db *gorm.DB) error {
	password := strings.TrimSpace(u.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedpassword)
	return nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

func GetUserByEmail(email string, db *gorm.DB) *User {

	user := &User{}
	err := db.Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil
	}
	return user
}

func (u *User) UpdateUser(db *gorm.DB) error {
	//err := db.Model(&u).Updates(map[string]interface{}{"first_name": u.FirstName, "last_name": u.LastName, "password": u.Password, "dummy": false}).Error
	err := db.Table("users").Where("email = ?", u.Email).Updates(User{FirstName: u.FirstName, LastName: u.LastName, Password: u.Password, Dummy: false}).Error
	return err
}
