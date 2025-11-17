package types

type Direction int8

const (
	XZminus Direction = iota - 4
	Zminus
	Yminus
	Xminus
	None
	X
	Y
	Z
	XZ
)
