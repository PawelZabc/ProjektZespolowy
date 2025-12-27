package config

import rl "github.com/gen2brain/raylib-go/raylib"

// Physics
const (
	Gravity = float32(0.005)
)

// Player
var (
	PlayerSpawnpoint = rl.NewVector3(0, 0, 0)
)

const (
	PlayerRadius       = float32(0.5)
	PlayerHeight       = float32(1.0)
	PlayerCameraHeight = float32(0.8)
	PlayerSpeed        = float32(0.1)
	JumpStrength       = float32(0.1)
)

// Other
const (
	DefaultPort        = 9000
	ServerIp           = "127.0.0.1"
	NetworkBufferSize  = 1024
	ClientTimeoutTicks = int64(200)
)
