package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	camera    rl.Camera
	rotationX float32
	rotationY float32
}

func NewCamera() *Camera {
	return &Camera{
		camera: rl.Camera{
			Position:   rl.NewVector3(0, 4.0, 4.0),   // inicial pos
			Target:     rl.NewVector3(0.0, 1.0, 0.0), // inicial target
			Up:         rl.NewVector3(0.0, 1.0, 0.0), // never changes
			Fovy:       config.CameraFov,
			Projection: rl.CameraPerspective,
		},
		rotationX: -math32.Pi / 2,
		rotationY: -math32.Pi / 2,
	}
}

// Updates camera rotation based on mouse movement
func (c *Camera) Update(centerX, centerY int, mouseLocked bool) {
	if !mouseLocked {
		return
	}

	deltaMouse := rl.GetMousePosition()

	c.rotationX += (deltaMouse.X - float32(centerX)) / 100 * config.CameraSensivity
	c.rotationY -= (deltaMouse.Y - float32(centerY)) / 100 * config.CameraSensivity

	// Clamp vertical rotation
	if c.rotationY > config.CameraLockMax {
		c.rotationY = config.CameraLockMax
	} else if c.rotationY < config.CameraLockMin {
		c.rotationY = config.CameraLockMin
	}

	rl.SetMousePosition(centerX, centerY)
}

// Updates camera position to follow player
func (c *Camera) UpdatePosition(playerPos rl.Vector3) {
	// Calculate target direction
	target := c.calculateTargetDirection()

	// Position camera at player position with offset
	c.camera.Position = rl.Vector3Add(playerPos, rl.NewVector3(0, config.PlayerCameraHeight, 0))
	c.camera.Target = rl.Vector3Add(target, c.camera.Position)
}

// Calculates the direction vector the camera is looking
func (c *Camera) calculateTargetDirection() rl.Vector3 {
	target := rl.Vector3{
		X: float32(math32.Sin(c.rotationY) * math32.Cos(c.rotationX)),
		Z: float32(math32.Sin(c.rotationY) * math32.Sin(c.rotationX)),
		Y: float32(math32.Cos(c.rotationY)),
	}
	return rl.Vector3Normalize(target)
}

func (c *Camera) GetCamera() rl.Camera {
	return c.camera
}

func (c *Camera) GetRotationX() float32 {
	return c.rotationX
}

func (c *Camera) GetRotationY() float32 {
	return c.rotationY
}

func (c *Camera) GetPlayerCameraRay() colliders.Ray {
	direction := c.calculateTargetDirection()

	return colliders.Ray{
		Origin:    c.camera.Position,
		Direction: direction,
	}
}
