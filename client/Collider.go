package main

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
	PushbackFrom(Collider)
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

func (c *CubeCollider) PushbackFrom(c2 Collider) {
	if cube, ok := c2.(*CubeCollider); ok {
		c.PushbackFromCube(cube)
	} else if cylinder, ok := c2.(*CylinderCollider); ok {
		c.PushbackFromCylinder(cylinder)
	}

}

func (c *CylinderCollider) PushbackFrom(c2 Collider) {
	if cylinder, ok := c2.(*CylinderCollider); ok {
		c.PushbackFromCylinder(cylinder)

		// println(distanceXZ, distanceY1, distanceY2)

	} else if cube, ok := c2.(*CubeCollider); ok {
		c.PushbackFromCube(cube)
	}

}

func (c *CubeCollider) PushbackFromCube(cube *CubeCollider) {
	x1 := c.Position.X - (cube.Position.X + cube.SizeX)
	x2 := cube.Position.X - (c.Position.X + c.SizeX)
	y1 := c.Position.Y - (cube.Position.Y + cube.SizeY)
	y2 := cube.Position.Y - (c.Position.Y + c.SizeY)
	z1 := c.Position.Z - (cube.Position.Z + cube.SizeZ)
	z2 := cube.Position.Z - (c.Position.Z + c.SizeZ)
	// println(x1, x2, y1, y2, z1, z2)
	if x1 <= 0 && x2 <= 0 && y1 <= 0 && y2 <= 0 && z1 <= 0 && z2 <= 0 {
		xmax := math.Max(x1, x2)
		ymax := math.Max(y1, y2)
		zmax := math.Max(z1, z2)

		if xmax > zmax {
			if xmax > ymax {
				if x1 > x2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(-x1, 0, 0))
					print("right ")
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(x2, 0, 0))
					print("left ")
				}
			} else {
				if y1 > y2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -y1, 0))
					print("up ")
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, y2, 0))
					print("down ")
				}
			}
		} else {
			if zmax > ymax {
				if z1 > z2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, 0, -z1))
					print("foward ")
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, 0, z2))
					print("backward ")
				}
			} else {
				if y1 > y2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -y1, 0))
					print("up ")
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, y2, 0))
					print("down ")
				}

			}

		}
	}
}

func (c *CubeCollider) PushbackFromCylinder(cylinder *CylinderCollider) {
	diffrence := rl.Vector2Subtract(
		rl.Vector2{X: math.Min(c.Position.X+c.SizeX, math.Max(c.Position.X, cylinder.Position.X)),
			Y: math.Min(c.Position.Z+c.SizeZ, math.Max(c.Position.Z, cylinder.Position.Z))},
		rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z})
	distanceXZ := rl.Vector2Length(diffrence) - (cylinder.Radius)
	distanceY1 := c.Position.Y - (cylinder.Position.Y + cylinder.Height)
	distanceY2 := cylinder.Position.Y - (c.Position.Y + c.SizeY)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {
			print("side ")
			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))

		} else if distanceY1 > distanceY2 {
			print("up ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
		} else {
			print("down ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
		}
	}

}

func (c *CylinderCollider) PushbackFromCube(cube *CubeCollider) {
	diffrence := rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
		rl.Vector2{X: math.Min(cube.Position.X+cube.SizeX, math.Max(cube.Position.X, c.Position.X)),
			Y: math.Min(cube.Position.Z+cube.SizeZ, math.Max(cube.Position.Z, c.Position.Z))})
	distanceXZ := rl.Vector2Length(diffrence) - (c.Radius)
	distanceY1 := c.Position.Y - (cube.Position.Y + cube.SizeY)
	distanceY2 := cube.Position.Y - (c.Position.Y + c.Height)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {
			print("side ")
			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))

		} else if distanceY1 > distanceY2 {
			print("up ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
		} else {
			print("down ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
		}
	}

}

func (c *CylinderCollider) PushbackFromCylinder(cylinder *CylinderCollider) {
	diffrence := rl.Vector2Subtract(rl.Vector2{X: c.Position.X, Y: c.Position.Z},
		rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z})
	distanceXZ := rl.Vector2Length(diffrence) - (c.Radius + cylinder.Radius)
	distanceY1 := c.Position.Y - (cylinder.Position.Y + cylinder.Height)
	distanceY2 := cylinder.Position.Y - (c.Position.Y + c.Height)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {
			print("side ")
			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))

		} else if distanceY1 > distanceY2 {
			print("up ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
		} else {
			print("down ")
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
		}

	}

}
