package fontMeshCreator;

type Line struct {
	maxLength float32
	spaceSize float32
	words []Word
	currentLineLength float32
}

func CreateLine(spaceWidth, fontSize, maxLength float32) Line {
    return Line{
		SpaceSize: spaceWidth * fontSize;
		MaxLength: maxLength;
        Words: make([]Word, 0)
	}
}

func (l *Line) AttemptToAddWord(word Word) bool {
    additionalLength := word.WordWidth
    if words.IsEmpty() {
        additionalLength += l.SpaceSize
    }
	if l.CurrentLineLength + additionalLength <= l.MaxLength {
		l.Words = append(words, word)
		l.CurrentLineLength += additionalLength
		return true
	} else {
		return false
	}
}
