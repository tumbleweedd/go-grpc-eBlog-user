package client

import (
	"context"
	"fmt"
	pb2 "github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostServiceClient struct {
	Client pb2.PostServiceClient
}

func InitPostServiceClient(url string) PostServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := PostServiceClient{
		Client: pb2.NewPostServiceClient(cc),
	}

	return c
}
func (c *PostServiceClient) GetAllPostsByUserId(userId int) (*pb2.GetAllPostsByUserIdResponse, error) {
	req := &pb2.GetAllPostsByUserIdRequest{
		UserId: int64(userId),
	}

	return c.Client.GetAllPostsByUserId(context.Background(), req)
}
