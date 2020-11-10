package game

type settings struct {
	cameraSensitivity float32
}

func defaultSettings() *settings {
	return &settings{
		cameraSensitivity: 3,
	}
}