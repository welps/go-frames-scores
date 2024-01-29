package sports

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client interface {
	GetMatches(ctx context.Context, gameType GameType) (ClientMatchResponse, error)
	GetLiveMatches(ctx context.Context, gameType GameType) (ClientMatchResponse, error)
}

type client struct {
	apiHost string
	apiKey  string
	resty   *resty.Client
}

func NewClient(resty *resty.Client, apiHost, apiKey string) (Client, error) {
	if apiHost == "" || apiKey == "" {
		return nil, fmt.Errorf("invalid api host or key")
	}

	return &client{
		resty:   resty,
		apiHost: apiHost,
		apiKey:  apiKey,
	}, nil
}

func (c *client) getSportsID(gameType GameType) int {
	switch gameType {
	case Basketball:
		return 3
	case Tennis:
		return 2
	default:
		return 0
	}
}

func (c *client) GetMatches(ctx context.Context, gameType GameType) (ClientMatchResponse, error) {
	sportsID := c.getSportsID(gameType)
	if sportsID == 0 {
		return ClientMatchResponse{}, fmt.Errorf("invalid game type")
	}
	// TODO:: Pagination
	url := fmt.Sprintf("%s/sports/%d/events", c.apiHost, c.getSportsID(gameType))

	response, err := c.resty.R().
		SetHeader("x-rapidapi-key", c.apiKey).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return ClientMatchResponse{}, err
	}

	status := response.StatusCode()
	if status != 200 {
		return ClientMatchResponse{}, fmt.Errorf(
			"failed to get scores - status code %d, response body: %s",
			status,
			response.Body(),
		)
	}

	var result ClientMatchResponse
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return ClientMatchResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return result, nil
}

func (c *client) GetLiveMatches(ctx context.Context, gameType GameType) (ClientMatchResponse, error) {
	sportsID := c.getSportsID(gameType)
	if sportsID == 0 {
		return ClientMatchResponse{}, fmt.Errorf("invalid game type")
	}
	url := fmt.Sprintf("%s/sports/%d/events/live", c.apiHost, c.getSportsID(gameType))

	response, err := c.resty.R().
		SetHeader("x-rapidapi-key", c.apiKey).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return ClientMatchResponse{}, err
	}

	status := response.StatusCode()
	if status != 200 {
		return ClientMatchResponse{}, fmt.Errorf(
			"failed to get scores - status code %d, response body: %s",
			status,
			response.Body(),
		)
	}

	var result ClientMatchResponse
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return ClientMatchResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return result, nil
}
