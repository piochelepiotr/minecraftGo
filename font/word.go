package font

type Word struct {
	Characters []Character
	Width      float32
	FontSize   float32
}

func CreateWord(fontSize float32) Word {
	return Word{
		FontSize:   fontSize,
		Characters: make([]Character, 0),
	}
}

func (w *Word) AddCharacter(character Character) {
	w.Characters = append(w.Characters, character)
	w.Width += character.XAdvance * w.FontSize
}
