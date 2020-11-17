package block

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

// Block is an id representing a type of Block
type Block uint8

type TextureID int

// NumberRowsTextures is the number number of rows on the texture image
const NumberRowsTextures int = 19

const (
	anvilBase TextureID = 0
	anvilTopDamaged0 TextureID = 1
	anvilTopDamaged1 TextureID = 2
	anvilTopDamaged2 TextureID = 3
	beacon TextureID = 4
	bedRock TextureID = 5
	beetrootsStage0 TextureID = 6
	beetrootsStage1 TextureID = 7
	beetrootsStage2 TextureID = 8
	beetrootsStage3 TextureID = 9
	boneBlockSide TextureID = 10
	boneBlockTop TextureID = 11
	bookshelf TextureID = 12
	brewingStand TextureID = 13
	brewingStandBase TextureID = 14
	brick TextureID = 15
	cactusBottom TextureID = 16
	cactusSide TextureID = 17
	cactusTop TextureID = 18
	cakeBottom TextureID = 19
	cakeInner TextureID = 20
	cakeSide TextureID = 21
	cakeTop TextureID = 22
	carrotsStage0 TextureID = 23
	carrotsStage1 TextureID = 24
	carrotsStage2 TextureID = 25
	carrotsStage3 TextureID = 26
	cauldronBottom TextureID = 27
	cauldronInner TextureID = 28
	cauldronSide TextureID = 29
	cauldronTop TextureID = 30
	chainCommandBlockBack TextureID = 31
	chainCommandBlockConditional TextureID = 32
	chainCommandBlockFront TextureID = 33
	chainCommandBlockSide TextureID = 34
	chorusFlower TextureID = 35
	chorusFlowerDead TextureID = 36
	chorusPlant TextureID = 37
	clay TextureID = 38
	coalBlock TextureID = 39
	coalOre TextureID = 40
	cobblestone TextureID = 41
	cobblestoneMossy TextureID = 42
	cocoaStage0 TextureID = 43
	cocoaStage1 TextureID = 44
	cocoaStage2 TextureID = 45
	commandBlockBack TextureID = 46
	commandBlockConditional TextureID = 47
	commandBlockFront TextureID = 48
	commandBlockSide TextureID = 49
	comparatorOff TextureID = 50
	comparatorOn TextureID = 51
	craftingTableFront TextureID = 52
	craftingTableSide TextureID = 53
	craftingTableTop TextureID = 54
	deadbush TextureID = 55
	diamondBlock TextureID = 56
	diamondOre TextureID = 57
	dirt TextureID = 58
	dirtPodzolTop TextureID = 59
	dispenserFrontHorizontal TextureID = 60
	dispenserFrontVertical TextureID = 61
	doorAcaciaLower TextureID = 62
	doorAcaciaUpper TextureID = 63
	doorBirchLower TextureID = 64
	doorBirchUpper TextureID = 65
	doorDarkOakLower TextureID = 66
	doorDarkOakUpper TextureID = 67
	doorIronLower TextureID = 68
	doorIronUpper TextureID = 69
	doorJungleLower TextureID = 70
	doorJungleUpper TextureID = 71
	doorSpruceLower TextureID = 72
	doorSpruceUpper TextureID = 73
	doorWoodLower TextureID = 74
	doorWoodUpper TextureID = 75
	doublePlantFernBottom TextureID = 76
	doublePlantFernTop TextureID = 77
	doublePlantGrassBottom TextureID = 78
	doublePlantGrassTop TextureID = 79
	doublePlantPaeoniaBottom TextureID = 80
	doublePlantPaeoniaTop TextureID = 81
	doublePlantRoseBottom TextureID = 82
	doublePlantRoseTop TextureID = 83
	doublePlantSunflowerBack TextureID = 84
	doublePlantSunflowerBottom TextureID = 85
	doublePlantSunflowerFront TextureID = 86
	doublePlantSunflowerTop TextureID = 87
	doublePlantSyringaBottom TextureID = 88
	doublePlantSyringaTop TextureID = 89
	dropperFrontHorizontal TextureID = 90
	dropperFrontVertical TextureID = 91
	emeraldBlock TextureID = 92
	emeraldOre TextureID = 93
	enchantingTableBottom TextureID = 94
	enchantingTableSide TextureID = 95
	enchantingTableTop TextureID = 96
	endBricks TextureID = 97
	endRod TextureID = 98
	endStone TextureID = 99
	endframeEye TextureID = 100
	endframeSide TextureID = 101
	endframeTop TextureID = 102
	farmlandDry TextureID = 103
	farmlandWet TextureID = 104
	fern TextureID = 105
	flowerAllium TextureID = 106
	flowerBlueOrchid TextureID = 107
	flowerDandelion TextureID = 108
	flowerHoustonia TextureID = 109
	flowerOxeyeDaisy TextureID = 110
	flowerPaeonia TextureID = 111
	flowerPot TextureID = 112
	flowerRose TextureID = 113
	flowerTulipOrange TextureID = 114
	flowerTulipPink TextureID = 115
	flowerTulipRed TextureID = 116
	flowerTulipWhite TextureID = 117
	frostedIce0 TextureID = 118
	frostedIce1 TextureID = 119
	frostedIce2 TextureID = 120
	frostedIce3 TextureID = 121
	furnaceFrontOff TextureID = 122
	furnaceFrontOn TextureID = 123
	furnaceSide TextureID = 124
	furnaceTop TextureID = 125
	glass TextureID = 126
	glassBlack TextureID = 127
	glassBlue TextureID = 128
	glassBrown TextureID = 129
	glassCyan TextureID = 130
	glassGray TextureID = 131
	glassGreen TextureID = 132
	glassLightBlue TextureID = 133
	glassLime TextureID = 134
	glassMagenta TextureID = 135
	glassOrange TextureID = 136
	glassPaneTop TextureID = 137
	glassPaneTopBlack TextureID = 138
	glassPaneTopBlue TextureID = 139
	glassPaneTopBrown TextureID = 140
	glassPaneTopCyan TextureID = 141
	glassPaneTopGray TextureID = 142
	glassPaneTopGreen TextureID = 143
	glassPaneTopLightBlue TextureID = 144
	glassPaneTopLime TextureID = 145
	glassPaneTopMagenta TextureID = 146
	glassPaneTopOrange TextureID = 147
	glassPaneTopPink TextureID = 148
	glassPaneTopPurple TextureID = 149
	glassPaneTopRed TextureID = 150
	glassPaneTopSilver TextureID = 151
	glassPaneTopWhite TextureID = 152
	glassPaneTopYellow TextureID = 153
	glassPink TextureID = 154
	glassPurple TextureID = 155
	glassRed TextureID = 156
	glassSilver TextureID = 157
	glassWhite TextureID = 158
	glassYellow TextureID = 159
	glowstone TextureID = 160
	goldBlock TextureID = 161
	goldOre TextureID = 162
	grassPathSide TextureID = 163
	grassPathTop TextureID = 164
	grassSide TextureID = 165
	grassTop TextureID = 166
	gravel TextureID = 167
	hardenedClay TextureID = 168
	hayBlockSide TextureID = 169
	hayBlockTop TextureID = 170
	hopperInside TextureID = 171
	hopperOutside TextureID = 172
	hopperTop TextureID = 173
	ice TextureID = 174
	icePacked TextureID = 175
	ironBars TextureID = 176
	ironBlock TextureID = 177
	ironOre TextureID = 178
	ironTrapdoor TextureID = 179
	itemframeBackground TextureID = 180
	jukeboxSide TextureID = 181
	jukeboxTop TextureID = 182
	ladder TextureID = 183
	lapisBlock TextureID = 184
	lapisOre TextureID = 185
	leavesAcacia TextureID = 186
	leavesBigOak TextureID = 187
	leavesBirch TextureID = 188
	leavesJungle TextureID = 189
	leavesOak TextureID = 190
	leavesSpruce TextureID = 191
	lever TextureID = 192
	logAcacia TextureID = 193
	logAcaciaTop TextureID = 194
	logBigOak TextureID = 195
	logBigOakTop TextureID = 196
	logBirch TextureID = 197
	logBirchTop TextureID = 198
	logJungle TextureID = 199
	logJungleTop TextureID = 200
	logOak TextureID = 201
	logOakTop TextureID = 202
	logSpruce TextureID = 203
	logSpruceTop TextureID = 204
	magma TextureID = 205
	melonSide TextureID = 206
	melonStemConnected TextureID = 207
	melonStemDisconnected TextureID = 208
	melonTop TextureID = 209
	mobSpawner TextureID = 210
	mushroomBlockInside TextureID = 211
	mushroomBlockSkinBrown TextureID = 212
	mushroomBlockSkinRed TextureID = 213
	mushroomBlockSkinStem TextureID = 214
	mushroomBrown TextureID = 215
	mushroomRed TextureID = 216
	myceliumSide TextureID = 217
	myceliumTop TextureID = 218
	netherBrick TextureID = 219
	netherWartBlock TextureID = 220
	netherWartStage0 TextureID = 221
	netherWartStage1 TextureID = 222
	netherWartStage2 TextureID = 223
	netherrack TextureID = 224
	noteblock TextureID = 225
	observerBack TextureID = 226
	observerBackLit TextureID = 227
	observerFront TextureID = 228
	observerSide TextureID = 229
	observerTop TextureID = 230
	obsidian TextureID = 231
	pistonBottom TextureID = 232
	pistonInner TextureID = 233
	pistonSide TextureID = 234
	pistonTopNormal TextureID = 235
	pistonTopSticky TextureID = 236
	planksAcacia TextureID = 237
	planksBigOak TextureID = 238
	planksBirch TextureID = 239
	planksJungle TextureID = 240
	planksOak TextureID = 241
	planksSpruce TextureID = 242
	potatoesStage0 TextureID = 243
	potatoesStage1 TextureID = 244
	potatoesStage2 TextureID = 245
	potatoesStage3 TextureID = 246
	prismarineBricks TextureID = 247
	prismarineDark TextureID = 248
	prismarineRough TextureID = 249
	pumpkinFaceOff TextureID = 250
	pumpkinFaceOn TextureID = 251
	pumpkinSide TextureID = 252
	pumpkinTop TextureID = 253
	purpurBlock TextureID = 254
	purpurPillar TextureID = 255
	purpurPillarTop TextureID = 256
	quartzBlockBottom TextureID = 257
	quartzBlockChiseled TextureID = 258
	quartzBlockChiseledTop TextureID = 259
	quartzBlockLines TextureID = 260
	quartzBlockLinesTop TextureID = 261
	quartzBlockSide TextureID = 262
	quartzBlockTop TextureID = 263
	quartzOre TextureID = 264
	railActivator TextureID = 265
	railActivatorPowered TextureID = 266
	railDetector TextureID = 267
	railDetectorPowered TextureID = 268
	railGolden TextureID = 269
	railGoldenPowered TextureID = 270
	railNormal TextureID = 271
	railNormalTurned TextureID = 272
	redNetherBrick TextureID = 273
	redSand TextureID = 274
	redSandstoneBottom TextureID = 275
	redSandstoneCarved TextureID = 276
	redSandstoneNormal TextureID = 277
	redSandstoneSmooth TextureID = 278
	redSandstoneTop TextureID = 279
	redstoneBlock TextureID = 280
	redstoneLampOff TextureID = 281
	redstoneLampOn TextureID = 282
	redstoneOre TextureID = 283
	redstoneTorchOff TextureID = 284
	redstoneTorchOn TextureID = 285
	repeaterOff TextureID = 286
	repeaterOn TextureID = 287
	repeatingCommandBlockBack TextureID = 288
	repeatingCommandBlockConditional TextureID = 289
	repeatingCommandBlockFront TextureID = 290
	repeatingCommandBlockSide TextureID = 291
	sand TextureID = 292
	sandstoneBottom TextureID = 293
	sandstoneCarved TextureID = 294
	sandstoneNormal TextureID = 295
	sandstoneSmooth TextureID = 296
	sandstoneTop TextureID = 297
	saplingAcacia TextureID = 298
	saplingBirch TextureID = 299
	saplingJungle TextureID = 300
	saplingOak TextureID = 301
	saplingRoofedOak TextureID = 302
	saplingSpruce TextureID = 303
	slime TextureID = 304
	snow TextureID = 305
	soulSand TextureID = 306
	sponge TextureID = 307
	spongeWet TextureID = 308
	stone TextureID = 309
	stoneAndesite TextureID = 310
	stoneAndesiteSmooth TextureID = 311
	stoneDiorite TextureID = 312
	stoneDioriteSmooth TextureID = 313
	stoneGranite TextureID = 314
	stoneGraniteSmooth TextureID = 315
	stoneSlabSide TextureID = 316
	stoneSlabTop TextureID = 317
	stonebrick TextureID = 318
	stonebrickCarved TextureID = 319
	stonebrickCracked TextureID = 320
	stonebrickMossy TextureID = 321
	tallgrass TextureID = 322
	tntBottom TextureID = 323
	tntSide TextureID = 324
	tntTop TextureID = 325
	torchOn TextureID = 326
	trapdoor TextureID = 327
	tripWireSource TextureID = 328
	vine TextureID = 329
	vine2 TextureID = 330
	vine3 TextureID = 331
	vine4 TextureID = 332
	waterlily TextureID = 333
	web TextureID = 334
	wheatStage0 TextureID = 335
	wheatStage1 TextureID = 336
	wheatStage2 TextureID = 337
	wheatStage3 TextureID = 338
	wheatStage4 TextureID = 339
	wheatStage5 TextureID = 340
	wheatStage6 TextureID = 341
	wheatStage7 TextureID = 342
)

const (
	Grass        Block = 0
	Stone        Block = 1
	Dirt         Block = 2
	TallGrass Block = 3
	BirchSapling Block = 4
	BedRock      Block = 17
	Sand         Block = 18
	Birch         Block = 21
	Gold         Block = 32
	Iron         Block = 33
	Coal         Block = 34
	BirchLeaves Block = 52
	Cactus       Block = 69
	Air          Block = 255
)

// blockFaces is the mapping from Block to TextureID
//it allows you to put different blocks on sides, top and bottom of blocks
var blockFaces = map[Block]map[Face]TextureID{
	Grass: {
		Side:   grassSide,
		Bottom: dirt,
		Top: grassTop,
	},
	Birch: {
		Side: logBirch,
		Top: logBirchTop,
		Bottom: logBirchTop,
	},
	Cactus: {
		Side:   cactusSide,
		Bottom: cactusBottom,
		Top: cactusTop,
	},
	Stone: {
		Side: stone,
	},
	Coal: {
		Side: coalOre,
	},
	Iron: {
		Side: ironOre,
	},
	Gold: {
		Side: goldOre,
	},
	BirchLeaves: {
		Side: leavesBirch,
	},
	Sand: {
		Side: sand,
	},
	BedRock: {
		Side: bedRock,
	},
	Dirt: {
		Side: dirt,
	},
	TallGrass: {
		Side: tallgrass,
	},
	BirchSapling: {
		Side: saplingBirch,
	},
}

var blockColors = map[TextureID]mgl32.Vec3{
	leavesBirch: {55.0 / 255.0, 97.0 / 255.0, 43.0 / 255.0},
	// grassTop: {71.0 / 255.0, 113.0 / 255.0, 53.0 / 255.0},
	grassTop: {121.0 / 255.0, 193.0 / 255.0, 81.0 / 255.0},
	tallgrass: {71.0 / 255.0, 113.0 / 255.0, 53.0 / 255.0},
}

var crossBlocks = map[Block]struct{}{
	TallGrass: {},
	BirchSapling: {},
}

var transparentBlocks = map[Block]struct{}{
	BirchLeaves: {},
	Air: {},
}

type Face byte

const (
	Top    Face = 0
	Side   Face = 1
	Bottom Face = 2
)

func (b Block) IsSolid() bool {
	return !(b == Air || b.IsCrossBlock())
}

func (b Block) IsTransparent() bool {
	_, ok := transparentBlocks[b]
	return ok  || b.IsCrossBlock()
}

func (b Block) IsCrossBlock() bool {
	_, ok := crossBlocks[b]
	return ok
}

func (b Block) GetSide(f Face) TextureID {
	if sides, ok := blockFaces[b]; ok {
		if side, ok := sides[f]; ok {
			return side
		}
		if side, ok := sides[Side]; ok {
			return side
		}
	}
	log.Fatal("Failed to get side for block", b)
	return -1
}

func (t TextureID) GetColor() mgl32.Vec3 {
	if color, ok := blockColors[t]; ok {
		return color
	}
	return mgl32.Vec3{1, 1, 1}
}
