package sports

import (
	"encoding/json"
	"strconv"
)

type ClientMatchResponse struct {
	Matches []ClientMatch `json:"data"`
}

type ClientMatch struct {
	ID                   int               `json:"id"`
	SportID              int               `json:"sport_id"`
	HomeTeamID           int               `json:"home_team_id"`
	AwayTeamID           int               `json:"away_team_id"`
	LeagueID             int               `json:"league_id"`
	ChallengeID          int               `json:"challenge_id"`
	SeasonID             int               `json:"season_id"`
	VenueID              int               `json:"venue_id"`
	RefereeID            interface{}       `json:"referee_id"`
	Slug                 string            `json:"slug"`
	Name                 string            `json:"name"`
	Status               string            `json:"status"`
	StatusMore           string            `json:"status_more"`
	TimeDetails          ClientTimeDetails `json:"time_details"`
	HomeTeam             ClientTeam        `json:"home_team"`
	AwayTeam             ClientTeam        `json:"away_team"`
	StartAt              string            `json:"start_at"`
	Priority             int               `json:"priority"`
	HomeScore            ClientScore       `json:"home_score"`
	AwayScore            ClientScore       `json:"away_score"`
	WinnerCode           int               `json:"winner_code"`
	AggregatedWinnerCode interface{}       `json:"aggregated_winner_code"`
	ResultOnly           bool              `json:"result_only"`
	Coverage             interface{}       `json:"coverage"`
	GroundType           interface{}       `json:"ground_type"`
	RoundNumber          interface{}       `json:"round_number"`
	SeriesCount          int               `json:"series_count"`
	MediasCount          interface{}       `json:"medias_count"`
	StatusLineup         int               `json:"status_lineup"`
	FirstSupply          interface{}       `json:"first_supply"`
	CardsCode            interface{}       `json:"cards_code"`
	EventDataChange      interface{}       `json:"event_data_change"`
	LastedPeriod         string            `json:"lasted_period"`
	DefaultPeriodCount   int               `json:"default_period_count"`
	Attendance           interface{}       `json:"attendance"`
	CupMatchOrder        interface{}       `json:"cup_match_order"`
	CupMatchInRound      interface{}       `json:"cup_match_in_round"`
	Periods              ClientPeriods     `json:"periods"`
	RoundInfo            interface{}       `json:"round_info"`
	PeriodsTime          ClientPeriodsTime `json:"periods_time"`
	MainOdds             ClientMainOdds    `json:"main_odds"`
	League               ClientLeague      `json:"league"`
	Challenge            ClientChallenge   `json:"challenge"`
	Season               ClientSeason      `json:"season"`
	Section              ClientSection     `json:"section"`
	Sport                ClientSport       `json:"sport"`
}

type ClientTimeDetails struct {
	Played                      int `json:"played"`
	PeriodLength                int `json:"periodLength"`
	OvertimeLength              int `json:"overtimeLength"`
	TotalPeriodCount            int `json:"totalPeriodCount"`
	CurrentPeriodStartTimestamp int `json:"currentPeriodStartTimestamp"`
}

type ClientTeam struct {
	ID               int                `json:"id"`
	SportID          int                `json:"sport_id"`
	CategoryID       int                `json:"category_id"`
	VenueID          int                `json:"venue_id"`
	ManagerID        int                `json:"manager_id"`
	Slug             string             `json:"slug"`
	Name             string             `json:"name"`
	HasLogo          bool               `json:"has_logo"`
	Logo             string             `json:"logo"`
	NameTranslations ClientTranslations `json:"name_translations"`
	NameShort        string             `json:"name_short"`
	NameFull         string             `json:"name_full"`
	NameCode         string             `json:"name_code"`
	HasSub           bool               `json:"has_sub"`
	Gender           string             `json:"gender"`
	IsNationality    bool               `json:"is_nationality"`
	CountryCode      string             `json:"country_code"`
	Country          string             `json:"country"`
	Flag             string             `json:"flag"`
	Foundation       interface{}        `json:"foundation"`
}

type ClientTranslations struct {
	En string `json:"en"`
	Ru string `json:"ru"`
	De string `json:"de"`
	Zh string `json:"zh"`
	El string `json:"el"`
	Nl string `json:"nl"`
	Pt string `json:"pt"`
}

type ClientScore map[string]StringOrInt
type ClientPeriods struct {
	Current  string `json:"current"`
	Period1  string `json:"period_1"`
	Period2  string `json:"period_2"`
	Period3  string `json:"period_3"`
	Period4  string `json:"period_4"`
	Period5  string `json:"period_5"`
	OverTime string `json:"over_time"`
}
type ClientPeriodsTime struct {
	CurrentTime  int `json:"current_time"`
	Period1Time  int `json:"period_1_time"`
	Period2Time  int `json:"period_2_time"`
	Period3Time  int `json:"period_3_time"`
	Period4Time  int `json:"period_4_time"`
	Period5Time  int `json:"period_5_time"`
	OverTimeTime int `json:"over_time_time"`
}

type ClientMainOdds struct {
	Outcome1 ClientOdds `json:"outcome_1"`
	Outcome2 ClientOdds `json:"outcome_2"`
}

type ClientOdds struct {
	Value  float64 `json:"value"`
	Change int     `json:"change"`
}

type ClientLeague struct {
	ID               int                `json:"id"`
	SportID          int                `json:"sport_id"`
	SectionID        int                `json:"section_id"`
	Slug             string             `json:"slug"`
	Name             string             `json:"name"`
	NameTranslations ClientTranslations `json:"name_translations"`
	HasLogo          bool               `json:"has_logo"`
	Logo             string             `json:"logo"`
}

type ClientChallenge struct {
	ID               int                `json:"id"`
	SportID          int                `json:"sport_id"`
	LeagueID         int                `json:"league_id"`
	Slug             string             `json:"slug"`
	Name             string             `json:"name"`
	NameTranslations ClientTranslations `json:"name_translations"`
	Order            int                `json:"order"`
	Priority         int                `json:"priority"`
}

type ClientSeason struct {
	ID        int    `json:"id"`
	LeagueID  int    `json:"league_id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	YearStart int    `json:"year_start"`
	YearEnd   int    `json:"year_end"`
}

type ClientSection struct {
	ID       int    `json:"id"`
	SportID  int    `json:"sport_id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Flag     string `json:"flag"`
}

type ClientSport struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type StringOrInt string

func (soi *StringOrInt) UnmarshalJSON(data []byte) error {
	var asInt int
	if err := json.Unmarshal(data, &asInt); err == nil {
		*soi = StringOrInt(strconv.Itoa(asInt))
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}

	*soi = StringOrInt(asString)
	return nil
}

func (soi *StringOrInt) String() string {
	return string(*soi)
}
