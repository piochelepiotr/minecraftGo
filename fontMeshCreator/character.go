package fontMeshCreator

type Character struct {
	Id int;
	XTextureCoord float32
	YTextureCoord float32
	XMaxTextureCoord float32
	YMaxTextureCoord float32
	XOffset float32
	YOffset float32
	SizeX float32
	SizeY float32
	XAdvance float32
}

func CreateCharacter(id int, xTextureCoord, yTextureCoord, xTexSize, yTexSize, xOffset, yOffset, sizeX, sizeY, xAdvance float32) Character {
    return Character{
        Id : id,
        XTextureCoord : xTextureCoord,
        YTextureCoord : yTextureCoord,
        XOffset : xOffset,
        YOffset : yOffset,
        SizeX : sizeX,
        SizeY : sizeY,
        XMaxTextureCoord : xTexSize + xTextureCoord,
        YMaxTextureCoord : yTexSize + yTextureCoord,
        XAdvance : xAdvance,
    }
}

