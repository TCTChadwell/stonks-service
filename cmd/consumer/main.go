package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pstonks "stonks-service/proto"
)

type StonkServiceClient struct {
	conn   *grpc.ClientConn
	client pstonks.StonksServiceClient
}

func jsonLog(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("error marshaling response:", err)
	}
	log.Printf("%s\n", b)
}

// NewPoiServiceClient created client for poi service
func NewStonkClient(url string) (*StonkServiceClient, error) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}
	return &StonkServiceClient{
		conn:   conn,
		client: pstonks.NewStonksServiceClient(conn),
	}, nil
}

func (c *StonkServiceClient) GetFavorites(ctx context.Context, userId string, opts ...grpc.CallOption) (*pstonks.GetFavoritesResponse, error) {

	req := &pstonks.GetFavoritesRequest{
		UserId: userId,
	}
	resp, err := c.client.GetFavorites(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *StonkServiceClient) CreateFavorite(ctx context.Context, userId, symbol string, opts ...grpc.CallOption) (*pstonks.CreateFavoritesResponse, error) {

	req := &pstonks.CreateFavoritesRequest{
		UserId: userId,
		Symbol: symbol,
	}
	resp, err := c.client.CreateFavorite(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *StonkServiceClient) DeleteFavorite(ctx context.Context, userId, symbol string, opts ...grpc.CallOption) (*pstonks.DeleteFavoritesResponse, error) {

	req := &pstonks.DeleteFavoritesRequest{
		UserId: userId,
		Symbol: symbol,
	}
	resp, err := c.client.DeleteFavorite(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

var (
	svcHost       string
	svcPort       string
	userId        string
	stockSymbol   string
	delAfterwards bool
)

func main() {
	{
		flag.StringVar(&svcHost, "host", "0.0.0.0", "service host")
		flag.StringVar(&svcPort, "port", "9001", "service host")
		flag.StringVar(&userId, "user", "ABC", "user-id")
		flag.StringVar(&stockSymbol, "stock", "SPY", "stock symbol")
		flag.BoolVar(&delAfterwards, "remove", false, "remove after deletion")
	}

	flag.Parse()

	url := fmt.Sprintf("%v:%v", svcHost, svcPort)
	client, err := NewStonkClient(url)
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}

	fmt.Println("creating favorite: ", stockSymbol)
	ctx := context.Background()
	createResp, err := client.CreateFavorite(ctx, userId, stockSymbol)
	if err != nil {
		log.Fatalf("Unable to create favorite: %v", err)
	}

	jsonLog(createResp)

	if delAfterwards {
		fmt.Println("deleting favorite: ", stockSymbol)
		_, err = client.DeleteFavorite(ctx, userId, stockSymbol)
		if err != nil {
			log.Fatalf("Unable to delete favorite: %v", err)
		}
	}

	fmt.Println("getting favorites for user: ", userId)

	getResp, err := client.GetFavorites(ctx, userId)
	if err != nil {
		log.Fatalf("Unable to get favorites: ", err)
	}
	jsonLog(getResp)
}
