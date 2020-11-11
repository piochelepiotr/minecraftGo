package state

type ID byte

const (
	Empty    ID = 0
	Game     ID = 1
	GameMenu ID = 2
	MainMenu ID = 3
)

// Switch is the action of switching to a new state (open a menu for eg)
type Switch struct {
	ID ID
	// WorldName is only used when switching to the game state from the main menu
	WorldName string
}