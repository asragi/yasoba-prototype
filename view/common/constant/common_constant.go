package constant

const sizeType = 2

var sizeMap = map[int]map[int]int{
	0: {
		0: 384,
		1: 288,
	},
	1: {
		0: 416,
		1: 312,
	},
	2: {
		0: 448,
		1: 336,
	},
}

var GameWidth = sizeMap[sizeType][0]
var GameWidthHalf = GameWidth / 2
var GameHeight = sizeMap[sizeType][1]
var GameHeightHalf = GameHeight / 2
