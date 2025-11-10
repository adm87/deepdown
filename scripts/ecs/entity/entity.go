package entity

// =========== Entity ===========

type Entity uint32

const (
	Null Entity = Entity(0)
)

func (e Entity) Equals(other Entity) bool {
	return e == other
}

func (e Entity) IsNull() bool {
	return e == Null
}

func (e Entity) ID() uint32 {
	return uint32(e)
}
