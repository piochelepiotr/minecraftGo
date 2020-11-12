package game

import "time"

type settings struct {
	cameraSensitivity float32
	doublePressDelay time.Duration
}

func defaultSettings() *settings {
	return &settings{
		cameraSensitivity: 3,
		doublePressDelay: time.Millisecond * 200,
	}
}