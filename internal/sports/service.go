package sports

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Service interface {
	GetMatches(ctx context.Context, gameType GameType) ([]Match, error)
	UpdateMatches(ctx context.Context) error
}

type service struct {
	cache  map[GameType][]Match
	mutex  *sync.RWMutex
	client Client
}

func NewService(client Client) Service {
	return &service{
		cache:  make(map[GameType][]Match),
		mutex:  &sync.RWMutex{},
		client: client,
	}
}

func (s *service) GetMatches(_ context.Context, gameType GameType) ([]Match, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	matches, ok := s.cache[gameType]
	if !ok {
		return nil, fmt.Errorf("no matches found for %s", gameType)
	}

	return matches, nil
}

func (s *service) UpdateMatches(ctx context.Context) error {
	err := s.updateMatches(ctx, Tennis, FormatTennisScore)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) updateMatches(ctx context.Context, gameType GameType, scoringFunc ScoringFunc) error {
	zap.S().Infof("Updating matches for %s", gameType)

	response, err := s.client.GetLiveMatches(ctx, gameType)
	if err != nil {
		return err
	}

	matches := make([]Match, 0, len(response.Matches))
	for _, match := range response.Matches {
		score, err := scoringFunc(match)
		if err != nil {
			zap.S().Errorw("unable to format tennis score", zap.Error(err))
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
	s.cache[gameType] = matches

	return nil
}
