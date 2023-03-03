package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/client"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pb"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/repository"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	lis, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", viper.GetString("port"))

	r := repository.NewRepository(db)
	postSvc := client.InitPostServiceClient(viper.GetString("post_svc_url"))
	s := service.NewService(r, postSvc)

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("pkg/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
