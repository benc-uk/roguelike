package graphics

import "image/color"

var FgColour color.RGBA = ColourWhite
var BgColour color.RGBA = ColourTrans

var ColourWhite = color.RGBA{255, 255, 255, 255}
var ColourBlack = color.RGBA{0, 0, 0, 255}
var ColourCursor = color.RGBA{60, 207, 180, 255}
var ColourTrans = color.RGBA{0, 0, 0, 0}

var ColourStatus = color.RGBA{0x10, 0x50, 0x10, 255}
var ColourStatusRed = color.RGBA{0x80, 0, 0, 255}
var ColourLog = color.RGBA{0, 0, 0x30, 0x30}
var ColourInv = color.RGBA{0x40, 0x40, 0x40, 255}
var ColourDialog = color.RGBA{0x50, 0x30, 0x10, 120}

var ColourTitle = color.RGBA{11, 45, 127, 255}
