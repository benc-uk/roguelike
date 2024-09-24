// nolint
package engine

type furniture struct {
	entityBase
}

func (f *furniture) Type() entityType {
	return entityTypeFurniture
}

func (f *furniture) BlocksLOS() bool {
	return true
}

func (f *furniture) BlocksMove() bool {
	return true
}
