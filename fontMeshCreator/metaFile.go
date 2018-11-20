package fontMeshCreator

import (
    "io/ioutil"
    "strconv"
    "strings"
)

const (
	PAD_TOP int = 0
	PAD_LEFT int = 1
	PAD_BOTTOM int = 2
	PAD_RIGHT int = 3
	DESIRED_PADDING int = 3
	SPLITTER string = " "
	NUMBER_SEPARATOR string = ","
)

type MetaFile struct {
	aspectRatio float32
	verticalPerPixelSize float32
	horizontalPerPixelSize float32
	spaceWidth float32
	padding []int
	paddingWidth int
	paddingHeight int
	metaData map[int]Character
	values = map[string]string
}

func CreateMetaFile(file string, aspectRatio float32) MetaFile {
    mf := MetaFile{
        aspectRatio: aspectRatio,
    }
    data := readFile(file)
	mf.loadPaddingData()
	mf.loadLineSizes()
    imageWidth := getValueOfVariable("scaleW")
	mf.loadCharacterData(imageWidth)
    return mf
}

func (mf *MetaFile) processNextLine(line string) {
	mf.values = make(map[string]string)
    for _, part := range strings.Split(line, SPLITTER) {
        valuePairs := strings.Split(part, "=")
		if len(valuePairs) == 2 {
			values[valuePairs[0]] = valuePairs[1]
		}
	}
}

func (mf *MetaFile) getValueOfVariable(variable string) {
    n, err :=  int(strconv.ParseInt(mf.values[variable], 10, 64))
    if err != nil {
        panic(err)
    }
    return n
}

func (mf *MetaFile)  getValuesOfVariable(variable string) []int {
    numbers := strings.Split(mf.values[variable], NUMBER_SEPARATOR)
	actualValues = make([]int, len(numbers))
    for i := 0; i < len(actualValues); i++ {
		actualValues[i] = strconv.ParseInt(numbers[i], 10, 64)
	}
	return actualValues
}

func readFile(file string) []byte {
    dat, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err)
    }
    return dat
}

func (mf *MetaFile) loadPaddingData() {
	mf.processNextLine()
	mf.padding = mf.getValuesOfVariable("padding")
	mf.paddingWidth = mf.padding[PAD_LEFT] + mf.padding[PAD_RIGHT];
	mf.paddingHeight = mf.padding[PAD_TOP] + mf.padding[PAD_BOTTOM];
}

func (mf *MetaFile) loadLineSizes() {
	mf.processNextLine();
    lineHeightPixels := getValueOfVariable("lineHeight") - mf.paddingHeight
	mf.verticalPerPixelSize = LINE_HEIGHT / float32(lineHeightPixels)
	mf.horizontalPerPixelSize = mf.verticalPerPixelSize / mf.aspectRatio
}

func (mf *MetaFile) loadCharacterData(imageWidth int) {
	mf.processNextLine()
	mf.processNextLine()
	for processNextLine() {
        c := mf.loadCharacter(imageWidth)
		if c != nil {
			mf.metaData[c.Id] = c
		}
	}
}

func (mf *MetaFile) loadCharacter(imageSize int) Character {
    id := mf.getValueOfVariable("id")
	if id == SPACE_ASCII {
		mf.spaceWidth = (mf.getValueOfVariable("xadvance") - mf.paddingWidth) * mf.horizontalPerPixelSize
		return nil
	}
    xTex := (float32(mf.getValueOfVariable("x")) + (mf.padding[PAD_LEFT] - DESIRED_PADDING)) / imageSize
    yTex := (float32(mf.getValueOfVariable("y")) + (mf.padding[PAD_TOP] - DESIRED_PADDING)) / imageSize
    width := mf.getValueOfVariable("width") - (mf.paddingWidth - (2 * DESIRED_PADDING))
    height := mf.getValueOfVariable("height") - ((mf.paddingHeight) - (2 * DESIRED_PADDING))
	quadWidth = width * mf.horizontalPerPixelSize
	quadHeight = height * mf.verticalPerPixelSize
    xTexSize := float32(width) / float32(imageSize)
    yTexSize := float32(height) / float32(imageSize)
    xOff := (mf.getValueOfVariable("xoffset") + mf.padding[PAD_LEFT] - DESIRED_PADDING) * mf.horizontalPerPixelSize;
    yOff := (mf.getValueOfVariable("yoffset") + (mf.padding[PAD_TOP] - DESIRED_PADDING)) * mf.verticalPerPixelSize;
    xAdvance := (mf.getValueOfVariable("xadvance") - mf.paddingWidth) * mf.horizontalPerPixelSize;
	return CreateCharacter(id, xTex, yTex, xTexSize, yTexSize, xOff, yOff, quadWidth, quadHeight, xAdvance)
}
