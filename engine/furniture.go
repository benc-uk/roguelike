// nolint
package engine

// ============================================================================
// Furniture entities are things like doors, barrels, tables etc
// ============================================================================

type furniture struct {
	entityBase
}

func (f furniture) Type() entityType {
	return entityTypeFurniture
}

func (f furniture) BlocksLOS() bool {
	return true
}

func (f furniture) BlocksMove() bool {
	return true
}

func (f furniture) String() string {
	return "furn_" + f.id + "_" + f.instanceID
}
