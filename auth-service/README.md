# Auth Service

The authentication gRPC service written in Golang

## Setup Database

Before we get started, let's create a MySQL database. Run the following sql script to create database:

```sql
CREATE USER 'mqttuser'@'localhost' IDENTIFIED BY 'password';
CREATE DATABASE mqtt_example CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON mqtt_example.* TO 'mqttuser'@'localhost';
FLUSH PRIVILEGES;
```

## ENV

Required `ENV`:

```properties
MYSQL_HOST=localhost:3306
MYSQL_DB=mqtt_example
MYSQL_USER=mqttuser
MYSQL_PASSWORD=password
AES_KEY=a very very very very secret key
```

Optional `ENV`:

```properties
DUMMY_FULL_NAME=Test User
DUMMY_EMAIL=test@mqtt.com
DUMMY_USERNAME=mqtt
DUMMY_PASSWORD=mqtt
DUMMY_ABOUT=I am a test user
DUMMY_CLIENT_KEY=mqtt-client
DUMMY_CLIENT_SECRET=secret
DUMMY_CLIENT_DESCRIPTION=The MQTT clients of test user
```

**NOTE:** The `AES_KEY` must be **32 bytes**

## Proto File

```protobuf
syntax = "proto3";

package pb;

// The authentication service definition.
service Auth {
  // Authentication mqtt client
  rpc AuthClient (AuthRequest) returns (AuthResponse) {}
}

// The request message for AuthClient
message AuthRequest {
  string clientKey = 1;
  string username = 2;
  string clientSecret = 3;
}

// The request message for AuthClient
message AuthResponse {
  string clientKey = 1;
  string username = 2;
  string code = 3; // return 200 when success
  string detail = 4; // a human-readable explanation
}
```

## Error Code Table

Below is the error code table response by server

Code | Description
---- | -----------
**2xx** | Success
200     | Success authentication
**4xx** | Fail authentication
400     | Fail authentication, invalid credentials
404     | Fail authentication, client not found
**5xx** | Fail authentication
500     | Fail authentication, internal server error

## Start Server

To run `auth-service` rename file `sample.env` to `.env` and then run the following command:

```bash
$ go run main.go --migrate --dummy # run this for the first time, it will automatically create tables, and create a dummy data
$ go run main.go # nex time just run this to start the server
```
