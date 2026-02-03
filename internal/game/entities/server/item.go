package server

type ItemType uint8

const (
	ItemRepair = iota
)

type Item struct {
	Type         ItemType
	EffectObject *EffectObject
	Id           uint8
}
