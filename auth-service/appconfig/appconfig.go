package appconfig

import (
	"os"
)

type (
	// MySQL config
	MySQL struct {
		Host     string
		Database string
		User     string
		Password string
	}
	// Dummy a dummy data
	Dummy struct {
		FullName     string
		Email        string
		Username     string
		Password     string
		About        string
		ClientKey    string
		ClientSecret string
		Description  string
	}
	// Crypto a struct store crypto params
	Crypto struct {
		// AESKey must be 32 bytes string
		AESKey string
		// AESKey32 a 32 bytes array converted
		AESKey32 [32]byte
	}
	configuration struct {
		MySQL  *MySQL
		Dummy  *Dummy
		Crypto *Crypto
	}
)

// Config an application configuration
var Config *configuration

// Load the application configuration
func Load(mysqlHost, mysqlDB, mysqlUser, mysqlPassword string) {
	mysql := &MySQL{
		Host:     mysqlHost,
		Database: mysqlDB,
		User:     mysqlUser,
		Password: mysqlPassword,
	}
	dummy := &Dummy{
		FullName:     getEnv("DUMMY_FULL_NAME", "Test User"),
		Email:        getEnv("DUMMY_EMAIL", "test@mqtt.com"),
		Username:     getEnv("DUMMY_USERNAME", "mqtt"),
		Password:     getEnv("DUMMY_PASSWORD", "mqtt"),
		About:        getEnv("DUMMY_ABOUT", "I am a test user"),
		ClientKey:    getEnv("DUMMY_CLIENT_KEY", "mqtt-client"),
		ClientSecret: getEnv("DUMMY_CLIENT_SECRET", "secret"),
		Description:  getEnv("DUMMY_CLIENT_DESCRIPTION", "The MQTT clients of test user"),
	}
	crypto := &Crypto{
		AESKey:   os.Getenv("AES_KEY"),
		AESKey32: get32BytesKey(os.Getenv("AES_KEY")),
	}

	Config = &configuration{
		MySQL:  mysql,
		Dummy:  dummy,
		Crypto: crypto,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func get32BytesKey(key string) [32]byte {
	var key32 [32]byte
	copy(key32[:], []byte(key)[0:32])
	return key32
}
