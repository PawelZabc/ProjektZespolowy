package shared

type Direction int8

const (
	DirXZminus Direction = iota - 4
	DirZminus
	DirYminus
	DirXminus
	DirNone
	DirX
	DirY
	DirZ
	DirXZ
)
