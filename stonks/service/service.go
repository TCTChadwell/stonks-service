package service

import (
	"context"
	"log"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	// "github.com/sirupsen/logrus"
	gCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	gStatus "google.golang.org/grpc/status"

	pstonks "stonks-service/proto"

	"stonks-service/stonks/api"
	"stonks-service/stonks/config"
	"stonks-service/stonks/db"
)

type StonksService struct {
	dataStore *redis.Client
	priceApi  *api.PriceApiClient
}

func NewStonksService(c config.StonkServiceConfig) (*StonksService, error) {

	redisClient, err := db.InitRedisClient(c.RedisHost, c.RedisPassword)
	if err != nil {
		log.Printf("Unable to connect to redis: %v", err)
		return nil, err
	}

	priceApiClient, err := api.NewPriceApiClient(c.TradingApiKey)
	if err != nil {
		log.Printf("Unable to connect to stock price API: %v", err)
		return nil, err
	}

	return &StonksService{
		dataStore: redisClient,
		priceApi:  priceApiClient,
	}, nil
}

func (s *StonksService) GetFavorites(ctx context.Context, req *pstonks.GetFavoritesRequest) (*pstonks.GetFavoritesResponse, error) {
	log := GetLog(ctx)
	reqId := uuid.New().String()
	ctx = metadata.AppendToOutgoingContext(ctx, "request-id", reqId)

	if req.GetUserId() == "" {
		return nil, gStatus.Errorf(
			gCodes.InvalidArgument,
			"missing user id",
		)
	}
	log.WithFields(logrus.Fields{
		"user-id": req.GetUserId(),
	}).Info("fetching stonk favorites")

	symbolIds, err := db.GetFavorites(ctx, s.dataStore, req.GetUserId())
	if err != nil {
		if err == db.ErrNotFound || len(symbolIds) < 1 {
			return nil, gStatus.Errorf(
				gCodes.NotFound,
				"no favorites found",
			)
		}
		return nil, gStatus.Errorf(
			gCodes.Internal,
			"issue retrieving favorites: %v",
			err,
		)
	}

	if len(symbolIds) < 1 {
		return &pstonks.GetFavoritesResponse{}, nil
	}

	// GET prices
	stonks, err := s.priceApi.GetPrices(symbolIds)
	if err != nil {
		return nil, gStatus.Errorf(
			gCodes.Internal,
			"Issue fetching prices: %v",
			err,
		)
	}

	var outputStonks []*pstonks.Stonk
	for _, s := range stonks {
		new := &pstonks.Stonk{
			Symbol:           s.Symbol,
			CurrentPrice:     s.Price,
			Type:             s.Source,
			FiftyTwoWeekHigh: s.High,
			FiftyTwoWeekLow:  s.Low,
		}
		outputStonks = append(outputStonks, new)
	}

	resp := &pstonks.GetFavoritesResponse{
		Stonks: outputStonks,
	}

	return resp, nil

}

func (s *StonksService) CreateFavorite(ctx context.Context, req *pstonks.CreateFavoritesRequest) (*pstonks.CreateFavoritesResponse, error) {

	reqId := uuid.New().String()
	ctx = metadata.AppendToOutgoingContext(ctx, "request-id", reqId)

	if (req.GetUserId() == "") || (req.GetSymbol() == "") {
		return nil, gStatus.Errorf(
			gCodes.InvalidArgument,
			"must provide both user_id & symbol",
		)
	}

	log.Printf("fetching stock price for: %v", req.GetSymbol())
	stonks, err := s.priceApi.GetPrices([]string{req.GetSymbol()})
	if err != nil {
		log.Printf("querying favorites")
		return nil, gStatus.Errorf(
			gCodes.NotFound,
			"Unable to fetch price information for symbol: %v",
			err,
		)
	}

	if len(stonks) != 1 {
		return nil, gStatus.Errorf(
			gCodes.NotFound,
			"Cannot return a single symbol's price data: %v",
			err,
		)
	}

	err = db.AddFavorite(ctx, s.dataStore, req.GetUserId(), req.GetSymbol())
	if err != nil {
		return nil, gStatus.Errorf(
			gCodes.Internal,
			"issue storing favorites: %v",
			err,
		)
	}

	var resp pstonks.Stonk
	resp.Symbol = stonks[0].Symbol
	resp.CurrentPrice = stonks[0].Price
	resp.FiftyTwoWeekHigh = stonks[0].High
	resp.FiftyTwoWeekLow = stonks[0].Low
	resp.Type = stonks[0].Source

	return &pstonks.CreateFavoritesResponse{
		Stonk: &resp,
	}, nil

}

func (s *StonksService) DeleteFavorite(ctx context.Context, req *pstonks.DeleteFavoritesRequest) (*pstonks.DeleteFavoritesResponse, error) {

	reqId := uuid.New().String()
	ctx = metadata.AppendToOutgoingContext(ctx, "request-id", reqId)

	if (req.GetUserId() == "") || (req.GetSymbol() == "") {
		return nil, gStatus.Errorf(
			gCodes.InvalidArgument,
			"must provide both user_id & symbol",
		)
	}

	err := db.RemoveFavorite(ctx, s.dataStore, req.GetUserId(), req.GetSymbol())
	if err != nil {
		return nil, gStatus.Errorf(
			gCodes.Internal,
			"unable to delete favorite",
		)
	}

	return &pstonks.DeleteFavoritesResponse{}, nil
}
