package event

//Goal scored by the team
type Goal struct {
	//ID of the game
	GameID string
	//ID of the team sored
	TeamScoredID string
	//ID of the team playing at home
	HomeTeamID string
	//ID of the guest team
	GuestTeamID string
	//How many points were scored
	Value int
}
