package model

import "github.com/jinzhu/gorm"
import "golang.org/x/crypto/bcrypt"

type (
	// User is a mqtt client model
	User struct {
		gorm.Model
		FullName       string `gorm:"size:45"`
		Email          string `gorm:"type:varchar(100);unique_index"`
		Username       string `gorm:"type:varchar(25);unique_index"`
		Password       string `sql:"-"`
		HashedPassword string
		About          string
		Clients        []Client
	}
)

// BeforeCreate a gorm callback
func (u *User) BeforeCreate() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hashedPassword)
	return nil
}

// AllUsers list all clients
func (db *DB) AllUsers(users *[]User) []error {
	errs := db.Find(&users).GetErrors()
	return errs
}

// GetUserByEmail get user by give email
func (db *DB) GetUserByEmail(email string, user *User) []error {
	errs := db.Where(&User{Email: email}).First(&user).GetErrors()
	return errs
}

// GetUserByUsername get user by give username
func (db *DB) GetUserByUsername(username string, user *User) []error {
	errs := db.Where(&User{Username: username}).First(&user).GetErrors()
	return errs
}
