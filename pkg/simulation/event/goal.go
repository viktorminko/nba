package event

type Goal struct {
	GameID       string
	TeamScoredID string
	HomeTeamID   string
	GuestTeamID  string
	Value        int
}
