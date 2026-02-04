package server

type ItemType uint8

const (
	ItemNone = iota
	ItemRepair
)

type Item struct {
	Type         ItemType
	EffectObject *EffectObject
	Id           uint8
}
