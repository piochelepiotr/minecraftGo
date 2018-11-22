package fontMeshCreator

type Line struct {
	MaxLength         float32
	SpaceSize         float32
	Words             []Word
	CurrentLineLength float32
}

func CreateLine(spaceWidth, fontSize, maxLength float32) Line {
	return Line{
		SpaceSize: spaceWidth * fontSize,
		MaxLength: maxLength,
		Words:     make([]Word, 0),
	}
}

func (l *Line) AttemptToAddWord(word Word) bool {
	additionalLength := word.Width
	if len(l.Words) == 0 {
		additionalLength += l.SpaceSize
	}
	if l.CurrentLineLength+additionalLength <= l.MaxLength {
		l.Words = append(l.Words, word)
		l.CurrentLineLength += additionalLength
		return true
	} else {
		return false
	}
}
