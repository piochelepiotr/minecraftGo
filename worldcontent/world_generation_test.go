package worldcontent

import "testing"

func TestRandom(t *testing.T) {
	NewGenerator(Config{Seed: 234})
}
