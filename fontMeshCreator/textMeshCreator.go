package fontMeshCreator

const (
    LINE_HEIGHT float32 = 0.03
	SPACE_ASCII int = 32;
)

type TextMeshCreator struct {
	metaData MetaFile
}

func CreateTextMeshCreator(metaFile string) TextMeshCreator {
    return TextMeshCreator{
        metaData: CreateMetaFile(metaFile)
    }
}

func createTextMesh(text GUIText) TextMeshData {
    lines := createStructure(text)
    data := createQuadVertices(text, lines)
	return data
}

func (tmc *TextMeshCreator) createStructure(text GUIText) []Line {
    lines := make([]Line, 0)
	currentLine := CreateLine(metaData.SpaceWidth, text.FontSize, text.MaxLineSize)
    currentWord := CreateWord(text.FontSize)
	for _, c := text.TextString {
		if c == SPACE_ASCII {
            added := currentLine.attemptToAddWord(currentWord)
			if !added {
				lines = append(lines, currentLine)
				currentLine = CreateLine(metaData.SpaceWidth, text.FontSize, text.MaxLineSize)
				currentLine.attemptToAddWord(currentWord)
			}
			currentWord = CreateWord(text.FontSize)
			continue
		}
        character := tmc.metaData.Character(c)
		currentWord.addCharacter(character)
	}
	tmc.completeStructure(lines, currentLine, currentWord, text)
	return lines
}

func (tmc *TextMeshCreator) completeStructure(lines *[]Line, currentLine Line, currentWord Word, text GUIText) {
    added := currentLine.attemptToAddWord(currentWord)
	if !added {
		*lines = append(*lines, currentLine)
		currentLine = CreateLine(metaData.SpaceWidth, text.FontSize, text.MaxLineSize)
		currentLine.attemptToAddWord(currentWord)
	}
	*lines = append(*lines, currentLine)
}

func (tmc *TextMeshCreator) createQuadVertices(text GUIText, lines []Line) TextMeshData {
	text.NumberOfLines = len(lines.size)
    curserX := float32(0)
    curserY := float32(0)
    vertices := make([]float32, 0)
    textureCoords := make([]float32, 0)
	for _, line := lines {
		if text.Centered {
			curserX = (line.MaxLength - line.LineLength) / 2.0;
		}
        for word := line.words {
			for letter := word.Characters {
				tmc.addVerticesForCharacter(curserX, curserY, letter, text.FontSize, vertices)
				tmc.addTexCoords(textureCoords, letter.XTextureCoord, letter.YTextureCoord, letter.XMaxTextureCoord, letter.YMaxTextureCoord)
				curserX += letter.XAdvance * text.FontSize
			}
			curserX += metaData.SpaceWidth * text.FontSize
		}
		curserX = 0
		curserY += LINE_HEIGHT * text.FontSize
    }
	return CreateTextMeshData(vertices, textureCoords)
}

func addVerticesForCharacter(curserX, curserY float32, character Character, fontSize float32, vertices *[]float32) {
    x := curserX + (character.XOffset * fontSize)
    y := curserY + (character.YOffset * fontSize)
    maxX := x + (character.X * fontSize)
    maxY := y + (character.Y * fontSize)
    properX := (2 * x) - 1
    properY := (-2 * y) + 1
    properMaxX := (2 * maxX) - 1
    properMaxY := (-2 * maxY) + 1
	addVertices(vertices, properX, properY, properMaxX, properMaxY)
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

