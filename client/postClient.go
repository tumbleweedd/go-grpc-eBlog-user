package client

import (
	"context"
	"fmt"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostServiceClient struct {
	Client pb.PostServiceClient
}

func InitPostServiceClient(url string) PostServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := PostServiceClient{
		Client: pb.NewPostServiceClient(cc),
	}

	return c
}
func (c *PostServiceClient) GetAllPostsByUserId(userId int) (*pb.GetAllPostsByUserIdResponse, error) {
	req := &pb.GetAllPostsByUserIdRequest{
		UserId: int64(userId),
	}

	return c.Client.GetAllPostsByUserId(context.Background(), req)
}
