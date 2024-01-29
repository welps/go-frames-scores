package sports

//go:generate stringer -type=GameType

type GameType int

const (
	Unknown GameType = iota
	Basketball
	Tennis
)
