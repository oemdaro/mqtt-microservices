package model

import (
	"encoding/base64"

	"github.com/gtank/cryptopasta"
	"github.com/jinzhu/gorm"
	"github.com/oemdaro/mqtt-microservices-example/auth-service/appconfig"
)

type (
	// Client is a mqtt client model
	Client struct {
		gorm.Model
		UserID       uint `gorm:"index"` // Foreign key (belongs to)
		ClientKey    string
		ClientSecret string
		Description  string
	}
)

// BeforeCreate a gorm callback
func (c *Client) BeforeCreate() error {
	key := appconfig.Config.Crypto.AESKey32
	encryptedSecret, err := cryptopasta.Encrypt([]byte(c.ClientSecret), &key)
	if err != nil {
		return err
	}

	c.ClientSecret = base64.StdEncoding.EncodeToString(encryptedSecret)
	return nil
}

// AfterFind invokes required after loading a record from the database.
func (c *Client) AfterFind() {
	key := appconfig.Config.Crypto.AESKey32
	decodedSecret, _ := base64.StdEncoding.DecodeString(c.ClientSecret)
	decryptedSecret, _ := cryptopasta.Decrypt(decodedSecret, &key)
	c.ClientSecret = string(decryptedSecret)
}

// AllClients list all clients
func (db *DB) AllClients(clients *[]Client) []error {
	errs := db.Find(&clients).GetErrors()
	return errs
}

// GetClientsByUsername get clients of an user by give username
func (db *DB) GetClientsByUsername(username string, clients *[]Client) []error {
	errs := db.Joins("JOIN `users` ON `clients`.user_id = `users`.id AND `users`.username = ?", username).Find(&clients).GetErrors()
	return errs
}
