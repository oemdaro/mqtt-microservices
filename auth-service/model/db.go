package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // nolint: golint
	"github.com/oemdaro/mqtt-microservices/auth-service/appconfig"
)

// Datastore interface for easily mock up
type Datastore interface {
	AllUsers(*[]User) []error
	GetUserByUsername(string, *User) []error
	AllClients(*[]Client) []error
	GetClientsByUsername(string, *[]Client) []error
}

// DB is an instant hold database connection
type DB struct {
	*gorm.DB
}

// NewDB return a new MySQL connection
func NewDB() (*DB, error) {
	config := appconfig.Config.MySQL
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Database)
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	// configure connection pool

	return &DB{db}, nil
}
