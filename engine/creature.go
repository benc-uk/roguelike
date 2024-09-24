package engine

type creature struct {
	entityBase
}

func (f *creature) Type() entityType {
	return entityTypeCreature
}

func (f *creature) BlocksLOS() bool {
	return false
}

func (f *creature) BlocksMove() bool {
	return true
}
