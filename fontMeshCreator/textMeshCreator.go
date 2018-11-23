package fontMeshCreator

import (
	"fmt"
)

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
	curserX := float32(0)
	curserY := float32(0)
	vertices := make([]float32, 0)
	textureCoords := make([]float32, 0)
	if text.VCenterText {
		textHeight := LINE_HEIGHT / 2 * text.FontSize * float32(len(lines))
		fmt.Println("Text height is: ", textHeight)
		fmt.Println("max text height is", text.MaxTextHeight)
		curserY = (text.MaxTextHeight - textHeight) / 2.0
		fmt.Println("cursor pos:", curserY)
		//curserY = 0
	}
	for _, line := range lines {
		if text.CenterText {
			curserX = (line.MaxLength - line.CurrentLineLength) / 2.0
		}
		for _, word := range line.Words {
			for _, letter := range word.Characters {
				addVerticesForCharacter(curserX, curserY, letter, text.FontSize, &vertices)
				addTexCoords(&textureCoords, letter.XTextureCoord, letter.YTextureCoord, letter.XMaxTextureCoord, letter.YMaxTextureCoord)
				curserX += letter.XAdvance * text.FontSize
			}
			curserX += tmc.metaData.spaceWidth * text.FontSize
		}
		curserX = 0
		curserY += LINE_HEIGHT * text.FontSize
	}
	return TextMeshData{
		VertexPositions: vertices,
		TextureCoords:   textureCoords,
	}
}

//func (tmc *TextMeshCreator) GetLineHeight(windowHeight int) float32 {
//	return tmc.metaData.verticalPerPixelSize
//}

func addVerticesForCharacter(curserX, curserY float32, character Character, fontSize float32, vertices *[]float32) {
	x := curserX + (character.XOffset * fontSize)
	y := curserY + (character.YOffset * fontSize)
	fmt.Println("y:", y)
	maxX := x + (character.SizeX * fontSize)
	maxY := y + (character.SizeY * fontSize)
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
