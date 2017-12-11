//go:generate protoc -I ../pb --go_out=plugins=grpc:../pb ../pb/authclient.proto

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oemdaro/mqtt-microservices-example/auth-service/appconfig"
	"github.com/oemdaro/mqtt-microservices-example/auth-service/model"
	"github.com/oemdaro/mqtt-microservices-example/pb"
	"google.golang.org/grpc"
)

var (
	// mysqlHost MySQL hostname
	mysqlHost = flag.String("mysql-host", os.Getenv("MYSQL_HOST"), "The MySQL host to connect to")
	// mysqlDB MySQL database name
	mysqlDB = flag.String("mysql-db", os.Getenv("MYSQL_DB"), "The MySQL database name")
	// mysqlUser MySQL username
	mysqlUser = flag.String("mysql-user", os.Getenv("MYSQL_USER"), "The MySQL username")
	// mysqlPassword MySQL password
	mysqlPassword = flag.String("mysql-password", os.Getenv("MYSQL_PASSWORD"), "The MySQL password")
	// port gRPC port number
	port = flag.Int("port", 50051, "The server port")
	// migrate the schema
	migrate = flag.Bool("migrate", false, "Auto migrate the schema")
	// dummy insert dummy data
	dummy = flag.Bool("dummy", false, "Insert dummy data")
	// signals we want to gracefully shutdown when it receives a SIGTERM or SIGINT
	signals = make(chan os.Signal, 1)
	done    = make(chan bool, 1)
)

// Server is used to implement model.AuthClient
type Server struct {
	db model.Datastore
}

// AuthClient authenticate MQTT client
func (s *Server) AuthClient(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	var clients []model.Client
	errs := s.db.GetClientsByUsername(in.Username, &clients)
	if errs != nil {
		for _, err := range errs {
			if err == gorm.ErrRecordNotFound {
				return &pb.AuthResponse{
					ClientKey: in.ClientKey,
					Username:  in.Username,
					Code:      "404",
					Detail:    "client not found",
				}, nil
			}
		}
		return &pb.AuthResponse{
			ClientKey: in.ClientKey,
			Username:  in.Username,
			Code:      "500",
			Detail:    "an unknown error occurred",
		}, errs[0]
	}

	for _, client := range clients {
		if client.ClientKey == in.ClientKey {
			if in.ClientSecret == client.ClientSecret {
				return &pb.AuthResponse{
					ClientKey: in.ClientKey,
					Username:  in.Username,
					Code:      "200",
					Detail:    "success authentication",
				}, nil
			}
		}
	}
	return &pb.AuthResponse{
		ClientKey: in.ClientKey,
		Username:  in.Username,
		Code:      "400",
		Detail:    "invalid credentials",
	}, nil
}

func main() {
	flag.Parse()
	if *mysqlHost == "" || *mysqlDB == "" || *mysqlUser == "" || *mysqlPassword == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Load configuration
	appconfig.Load(*mysqlHost, *mysqlDB, *mysqlUser, *mysqlPassword)

	db, err := model.NewDB()
	if err != nil {
		log.Fatalf("Error connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err == nil {
			log.Println("Shut down completed")
		}
	}()

	if *migrate {
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Client{})
		log.Println("Done migrate database")
	}
	if *dummy {
		dummy := appconfig.Config.Dummy
		user := model.User{
			FullName: dummy.FullName,
			Email:    dummy.Email,
			Username: dummy.Username,
			Password: dummy.Password,
			About:    dummy.About,
			Clients: []model.Client{
				{ClientKey: dummy.ClientKey + "-1", ClientSecret: dummy.ClientSecret, Description: dummy.Description},
				{ClientKey: dummy.ClientKey + "-2", ClientSecret: dummy.ClientSecret, Description: dummy.Description},
			},
		}
		db.Create(&user)
		log.Println("Done create dummy data")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &Server{db})

	// Notify when receive SIGINT or SIGTERM
	// kill -SIGINT <PID> or Ctrl+c
	// kill -SIGTERM <PID>
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-signals:
				log.Println("Graceful shutting down...")
				log.Println("Stopping qRPC server...")
				s.GracefulStop()
				time.Sleep(time.Second)
				done <- true
			}
		}
	}()

	// Serve gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Exiting
	<-done
	log.Println("Closing database connection...")
}
