package main

import "github.com/hajimehoshi/ebiten/v2"

var controls = map[string][]ebiten.Key{
	"up":    {ebiten.KeyW, ebiten.KeyUp},
	"down":  {ebiten.KeyS, ebiten.KeyDown},
	"left":  {ebiten.KeyA, ebiten.KeyLeft},
	"right": {ebiten.KeyD, ebiten.KeyRight},
}
