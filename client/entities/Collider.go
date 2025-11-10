package entities

import (
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// type Vector = rl.Vector3

// func (v *Vector) Add(vec rl.Vector3) {
// 	v.X += vec.X
// }

type Collider interface {
	CollidesWith(Collider) bool
	GetPosition() rl.Vector3
	SetPosition(rl.Vector3)
	AddPosition(rl.Vector3)
}

type CylinderCollider struct {
	Position rl.Vector3
	Radius   float32
	Height   float32
}

func (c CylinderCollider) CollidesWith(c2 Collider) bool {
	if cylinder, ok := c2.(*CylinderCollider); ok {
		if rl.Vector2Distance(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
			rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z}) < (c.Radius+cylinder.Radius) &&
			c.Position.Y <= cylinder.Position.Y+cylinder.Height && c.Position.Y+c.Height >= cylinder.Position.Y {
			return true
		}

	} else if cube, ok := c2.(*CubeCollider); ok {
		if rl.Vector2Distance(rl.Vector2{X: math.Min(cube.Position.X+cube.SizeX, math.Max(cube.Position.X, c.Position.X)),
			Y: math.Min(cube.Position.Z+cube.SizeZ, math.Max(cube.Position.Z, c.Position.Z))},
			rl.Vector2{X: c.Position.X, Y: c.Position.Z}) <= c.Radius &&
			cube.Position.Y <= c.Position.Y+c.Height && cube.Position.Y+cube.SizeY >= c.Position.Y {
			return true
		}
	}

	return false
}

func (c CylinderCollider) GetPosition() rl.Vector3 {
	return c.Position
}

func (c *CylinderCollider) SetPosition(vec rl.Vector3) {
	c.Position = vec
}

func (c *CylinderCollider) AddPosition(vec rl.Vector3) {
	c.Position = rl.Vector3Add(c.Position, vec)
}

type CubeCollider struct {
	Position rl.Vector3
	SizeX    float32
	SizeY    float32
	SizeZ    float32
}

func (c CubeCollider) CollidesWith(c2 Collider) bool {
	if cylinder, ok := c2.(*CylinderCollider); ok {
		if rl.Vector2Distance(rl.Vector2{X: math.Min(c.Position.X+c.SizeX, math.Max(c.Position.X, cylinder.Position.X)),
			Y: math.Min(c.Position.Z+c.SizeZ, math.Max(c.Position.Z, cylinder.Position.Z))},
			rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z}) < cylinder.Radius &&
			c.Position.Y <= cylinder.Position.Y+cylinder.Height && c.Position.Y+c.SizeY >= cylinder.Position.Y {

			return true
		}
	} else if cube, ok := c2.(*CubeCollider); ok {
		if c.Position.Y <= cube.Position.Y+cube.SizeY && c.Position.Y+c.SizeY >= cube.Position.Y &&
			c.Position.X <= cube.Position.X+cube.SizeX && c.Position.X+c.SizeX >= cube.Position.X &&
			c.Position.Z <= cube.Position.Z+cube.SizeZ && c.Position.Z+c.SizeZ >= cube.Position.Z {
			return true
		}
	}

	return false
}

func (c CubeCollider) GetPosition() rl.Vector3 {
	return c.Position
}

func (c *CubeCollider) SetPosition(vec rl.Vector3) {
	c.Position = vec
}

func (c *CubeCollider) AddPosition(vec rl.Vector3) {
	c.Position = rl.Vector3Add(c.Position, vec)
}
