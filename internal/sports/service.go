package sports

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Service interface {
	GetMatches(ctx context.Context, gameType GameType, live bool) ([]Match, error)
	UpdateMatches(ctx context.Context, live bool) error
}

type service struct {
	cache  map[string][]Match
	mutex  *sync.RWMutex
	client Client
}

func NewService(client Client) Service {
	return &service{
		cache:  make(map[string][]Match),
		mutex:  &sync.RWMutex{},
		client: client,
	}
}

func (s *service) GetMatches(_ context.Context, gameType GameType, live bool) ([]Match, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	matches, ok := s.cache[s.getKey(gameType, live)]
	if !ok {
		return nil, fmt.Errorf("no matches found for %s", gameType)
	}

	return matches, nil
}

func (s *service) UpdateMatches(ctx context.Context, live bool) error {
	err := s.updateMatches(ctx, Tennis, FormatTennisScore, live)
	if err != nil {
		return err
	}
	err = s.updateMatches(ctx, Basketball, FormatBasketballScore, live)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) updateMatches(ctx context.Context, gameType GameType, scoringFunc ScoringFunc, live bool) error {
	zap.S().Infof("Updating matches for %s", gameType)

	var response ClientMatchResponse
	var err error
	if live {
		response, err = s.client.GetLiveMatches(ctx, gameType)
	} else {
		response, err = s.client.GetMatches(ctx, gameType)
	}
	if err != nil {
		return fmt.Errorf("unable to get matches: %w", err)
	}

	matches := make([]Match, 0, len(response.Matches))
	for _, match := range response.Matches {
		score, err := scoringFunc(match)
		if err != nil {
			zap.S().Errorw(fmt.Sprintf("unable to format %s score", gameType), zap.Error(err))
			continue
		}

		matches = append(
			matches, Match{
				GameType: gameType,
				Home:     Team{Name: match.HomeTeam.Name},
				Away:     Team{Name: match.AwayTeam.Name},
				Score:    score,
			},
		)
	}

	zap.S().Infof("Updated %d %s matches", len(matches), gameType)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache[s.getKey(gameType, live)] = matches

	return nil
}

func (s *service) getKey(gameType GameType, live bool) string {
	var liveStr string
	if live {
		liveStr = "live"
	}

	return fmt.Sprintf("%s_%s", gameType, liveStr)
}
