package sports

type Team struct {
	Name string
}
type Match struct {
	GameType GameType
	Home     Team
	Away     Team
	Score    Score
}
