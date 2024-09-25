package engine

import "roguelike/core"

// We're redeclaring some core types in the engine package, but why?
// Note they are lowercase, so they are not exported
// This means we can embed them in our structs _without_ exporting them from the engine package

type pos = core.Pos
type rect = core.Rect
type size = core.Size
type direction = core.Direction
