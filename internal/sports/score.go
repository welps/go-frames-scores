package sports

import (
	"fmt"
	"strconv"
	"strings"
)

type ScoringFunc func(ClientMatch) (Score, error)

type Score struct {
	Home      []string
	HomeTotal string
	Away      []string
	AwayTotal string
}

func FormatBasketballScore(match ClientMatch) (Score, error) {
	lp := strings.TrimPrefix(match.LastedPeriod, "period_")
	lastPeriod, err := strconv.Atoi(lp)
	if err != nil {
		return Score{}, fmt.Errorf("unable to determine last period: %w", err)
	}

	score := Score{
		Home: make([]string, 0, lastPeriod),
		Away: make([]string, 0, lastPeriod),
	}

	// TODO:: Needs a fix for overtime
	periodKeys := getPeriodKeys(lastPeriod)

	for _, period := range periodKeys {
		homeScore, err := getScoreForPeriod(match.HomeScore, period)
		if err != nil {
			return Score{}, fmt.Errorf("unable to get home score for period %s: %w", period, err)
		}
		score.Home = append(score.Home, homeScore)

		awayScore, err := getScoreForPeriod(match.AwayScore, period)
		if err != nil {
			return Score{}, fmt.Errorf("unable to get away score for period %s: %w", period, err)
		}
		score.Away = append(score.Away, awayScore)
	}

	return score, nil
}

func FormatTennisScore(match ClientMatch) (Score, error) {
	lp := strings.TrimPrefix(match.LastedPeriod, "period_")
	lastPeriod, err := strconv.Atoi(lp)
	if err != nil {
		return Score{}, fmt.Errorf("unable to determine last period: %w", err)
	}

	score := Score{
		Home: make([]string, 0, lastPeriod),
		Away: make([]string, 0, lastPeriod),
	}

	periodKeys := getPeriodKeys(lastPeriod)

	for _, period := range periodKeys {
		homeScore, err := getScoreForPeriod(match.HomeScore, period)
		if err != nil {
			return Score{}, fmt.Errorf("unable to get home score for period %s: %w", period, err)
		}
		score.Home = append(score.Home, homeScore)

		awayScore, err := getScoreForPeriod(match.AwayScore, period)
		if err != nil {
			return Score{}, fmt.Errorf("unable to get away score for period %s: %w", period, err)
		}
		score.Away = append(score.Away, awayScore)
	}

	return score, nil
}

func getScoreForPeriod(score ClientScore, period string) (string, error) {
	homeScore, ok := score[period]
	if !ok {
		return "", fmt.Errorf("unable to find score for period %s", period)
	}
	return homeScore.String(), nil
}

func getPeriodKeys(lastPeriod int) []string {
	periods := make([]string, 0, lastPeriod+1)
	for i := 0; i < lastPeriod; i++ {
		periods = append(periods, fmt.Sprintf("period_%d", i+1))
	}

	// TODO:: What is this for?
	//periods = append(periods, "current")

	return periods
}
