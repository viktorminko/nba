package game

//Game represents game happening between 2 teams
//one team plays at home another is a guest
type Game struct {
	ID          string
	Home, Guest *Team
}
