package protocol

import (
	"bytes"
	"encoding/binary"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/input"
)

type ClientData struct {
	Inputs    []input.PlayerAction
	RotationX float32
	RotationY float32
}

func SerializeClientData(c ClientData) []byte {
	size := 8 + 8 + len(c.Inputs)
	buf := make([]byte, 0, size)
	b := bytes.NewBuffer(buf)

	binary.Write(b, binary.LittleEndian, uint8(len(c.Inputs)))
	binary.Write(b, binary.LittleEndian, c.RotationX)
	binary.Write(b, binary.LittleEndian, c.RotationY)
	for _, input := range c.Inputs {
		binary.Write(b, binary.LittleEndian, input)
	}

	return b.Bytes()
}

func DeserializeClientData(buf []byte) ClientData {
	b := bytes.NewReader(buf)
	c := ClientData{}
	var size uint8
	binary.Read(b, binary.LittleEndian, &size)
	binary.Read(b, binary.LittleEndian, &c.RotationX)
	binary.Read(b, binary.LittleEndian, &c.RotationY)
	c.Inputs = make([]input.PlayerAction, size)
	for i := range size {
		binary.Read(b, binary.LittleEndian, &c.Inputs[i])
	}
	return c

}
