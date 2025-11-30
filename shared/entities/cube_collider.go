package entities

import (
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	math "github.com/chewxy/math32"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CubeCollider struct {
	Position rl.Vector3
	SizeX    float32
	SizeY    float32
	SizeZ    float32
}

func (c CubeCollider) CollidesWith(c2 types.Collider) bool {
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

func (c *CubeCollider) PushbackFrom(c2 types.Collider) types.Direction {
	if cube, ok := c2.(*CubeCollider); ok {
		return c.PushbackFromCube(cube)
	} else if cylinder, ok := c2.(*CylinderCollider); ok {
		return c.PushbackFromCylinder(cylinder)
	}

	return types.DirNone

}

func (c *CubeCollider) PushbackFromCube(cube *CubeCollider) types.Direction {
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
					return types.DirX
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(x2, 0, 0))
					return -types.DirX
				}
			} else {
				if y1 > y2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -y1, 0))
					return -types.DirY
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, y2, 0))
					return types.DirY
				}
			}
		} else {
			if zmax > ymax {
				if z1 > z2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, 0, -z1))
					return types.DirZ
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, 0, z2))
					return -types.DirZ
				}
			} else {
				if y1 > y2 {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -y1, 0))
					return -types.DirY
				} else {
					c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, y2, 0))
					return types.DirY
				}

			}

		}
	} else {
		return types.DirNone
	}
}

func (c *CubeCollider) PushbackFromCylinder(cylinder *CylinderCollider) types.Direction {
	diffrence := rl.Vector2Subtract(
		rl.Vector2{X: math.Min(c.Position.X+c.SizeX, math.Max(c.Position.X, cylinder.Position.X)),
			Y: math.Min(c.Position.Z+c.SizeZ, math.Max(c.Position.Z, cylinder.Position.Z))},
		rl.Vector2{X: cylinder.Position.X, Y: cylinder.Position.Z})
	distanceXZ := rl.Vector2Length(diffrence) - (cylinder.Radius)
	distanceY1 := c.Position.Y - (cylinder.Position.Y + cylinder.Height)
	distanceY2 := cylinder.Position.Y - (c.Position.Y + c.SizeY)

	if distanceXZ <= 0 && distanceY1 <= 0 && distanceY2 <= 0 {
		if distanceXZ > distanceY1 && distanceXZ > distanceY2 {

			forceXZ := rl.Vector2Scale(rl.Vector2Normalize(diffrence), -distanceXZ)
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(forceXZ.X, 0, forceXZ.Y))
			return types.DirXZ

		} else if distanceY1 > distanceY2 {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, -distanceY1, 0))
			return -types.DirY
		} else {
			c.Position = rl.Vector3Add(c.Position, rl.NewVector3(0, distanceY2, 0))
			return types.DirY
		}
	} else {
		return types.DirNone
	}

}
