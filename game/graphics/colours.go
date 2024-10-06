package graphics

import "image/color"

var FgColour color.RGBA = ColourWhite
var BgColour color.RGBA = ColourTrans

var ColourWhite = color.RGBA{0xff, 0xff, 0xff, 0xff}
var ColourBlack = color.RGBA{0x00, 0x00, 0x00, 0xff}
var ColourTurq = color.RGBA{0x00, 0xff, 0xff, 0xff}
var ColourTrans = color.RGBA{0, 0, 0, 0}

var ColourStatus = color.RGBA{0x10, 0x50, 0x10, 0xff}
var ColourStatusRed = color.RGBA{0x80, 0x00, 0x00, 0xff}
var ColourLog = color.RGBA{0x00, 0x00, 0x30, 0x30}
var ColourInv = color.RGBA{0x40, 0x40, 0x40, 0xff}
var ColourDialog = color.RGBA{0x50, 0x30, 0x10, 120}

var ColourTitle = color.RGBA{11, 45, 127, 255}
