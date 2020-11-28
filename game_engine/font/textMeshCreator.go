package font

const (
	LINE_HEIGHT float32 = 0.05
	SPACE_ASCII int     = 32
)

type TextMeshCreator struct {
	metaData MetaFile
}

func CreateTextMeshCreator(metaFile string, aspectRatio float32) TextMeshCreator {
	return TextMeshCreator{
		metaData: CreateMetaFile(metaFile, aspectRatio),
	}
}

func (tmc *TextMeshCreator) createTextMesh(text GUIText) TextMeshData {
	lines := tmc.createStructure(text)
	data := tmc.createQuadVertices(text, lines)
	return data
}

func (tmc *TextMeshCreator) createStructure(text GUIText) []Line {
	lines := make([]Line, 0)
	currentLine := CreateLine(tmc.metaData.spaceWidth, text.FontSize, text.MaxLineSize)
	currentWord := CreateWord(text.FontSize)
	for _, c := range text.TextString {
		char := int(c)
		if char == SPACE_ASCII {
			added := currentLine.AttemptToAddWord(currentWord)
			if !added {
				lines = append(lines, currentLine)
				currentLine = CreateLine(tmc.metaData.spaceWidth, text.FontSize, text.MaxLineSize)
				currentLine.AttemptToAddWord(currentWord)
			}
			currentWord = CreateWord(text.FontSize)
			continue
		}
		character := tmc.metaData.getCharacter(char)
		currentWord.AddCharacter(character)
	}
	tmc.completeStructure(&lines, currentLine, currentWord, text)
	return lines
}

func (tmc *TextMeshCreator) completeStructure(lines *[]Line, currentLine Line, currentWord Word, text GUIText) {
	added := currentLine.AttemptToAddWord(currentWord)
	if !added {
		*lines = append(*lines, currentLine)
		currentLine = CreateLine(tmc.metaData.spaceWidth, text.FontSize, text.MaxLineSize)
		currentLine.AttemptToAddWord(currentWord)
	}
	*lines = append(*lines, currentLine)
}

func (tmc *TextMeshCreator) createQuadVertices(text GUIText, lines []Line) TextMeshData {
	text.NumberOfLines = len(lines)
	cursorXStart := float32(0)
	cursorY := float32(0)
	vertices := make([]float32, 0)
	textureCoords := make([]float32, 0)
	for _, line := range lines {
		cursorX := cursorXStart
		for _, word := range line.Words {
			for _, letter := range word.Characters {
				addVerticesForCharacter(cursorX, cursorY, letter, text.FontSize, &vertices)
				addTexCoords(&textureCoords, letter.XTextureCoord, letter.YTextureCoord, letter.XMaxTextureCoord, letter.YMaxTextureCoord)
				cursorX += letter.XAdvance * text.FontSize
			}
			cursorX += tmc.metaData.spaceWidth * text.FontSize
		}
		if cursorX - cursorXStart > text.Width {
			text.Width = cursorX - cursorXStart
		}
		cursorY += LINE_HEIGHT * text.FontSize
	}
	return TextMeshData{
		VertexPositions: vertices,
		TextureCoords:   textureCoords,
		Width: text.Width,
	}
}

func pos(x float32) float32 {
	return 2*x
}

func addVerticesForCharacter(cursorX, cursorY float32, character Character, fontSize float32, vertices *[]float32) {
	x := pos(cursorX + (character.XOffset * fontSize))
	y := pos(cursorY + (character.YOffset * fontSize))
	maxX := x + pos(character.SizeX * fontSize)
	maxY := y + pos(character.SizeY * fontSize)
	// properY := (-2 * y) + 1
	// properMaxY := (-2 * maxY) + 1
	addVertices(vertices, x, y, maxX, maxY)
	// addVertices(vertices, x, y-1, maxX, maxY-1)
}

func addVertices(vertices *[]float32, x, y, maxX, maxY float32) {
	*vertices = append(*vertices, x)
	*vertices = append(*vertices, y)
	*vertices = append(*vertices, x)
	*vertices = append(*vertices, maxY)
	*vertices = append(*vertices, maxX)
	*vertices = append(*vertices, maxY)
	*vertices = append(*vertices, maxX)
	*vertices = append(*vertices, maxY)
	*vertices = append(*vertices, maxX)
	*vertices = append(*vertices, y)
	*vertices = append(*vertices, x)
	*vertices = append(*vertices, y)
}

func addTexCoords(texCoords *[]float32, x, y, maxX, maxY float32) {
	*texCoords = append(*texCoords, x)
	*texCoords = append(*texCoords, y)
	*texCoords = append(*texCoords, x)
	*texCoords = append(*texCoords, maxY)
	*texCoords = append(*texCoords, maxX)
	*texCoords = append(*texCoords, maxY)
	*texCoords = append(*texCoords, maxX)
	*texCoords = append(*texCoords, maxY)
	*texCoords = append(*texCoords, maxX)
	*texCoords = append(*texCoords, y)
	*texCoords = append(*texCoords, x)
	*texCoords = append(*texCoords, y)
}
