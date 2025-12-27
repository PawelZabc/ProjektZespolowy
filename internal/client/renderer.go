package client

import (
	"strconv"

	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Go there if you want to draw things in game
type Renderer struct {
	debugMode bool // not so future debug mode
}

func NewRenderer(debugMode bool) *Renderer {
	return &Renderer{
		debugMode: debugMode,
	}
}

func (r *Renderer) RenderWorld(state *GameState) {
	// rendering 3d world
	levels.DrawRoom(state.GetCurrentRoom())
	entities.DrawActorsMap(state.GetPlayers())
	state.GetEnemy().Draw()
	r.renderLights(state.lights)

	if r.debugMode {
		r.renderDebug(state)
	}
}

// randering 2d elements
func (r *Renderer) RenderUI(state *GameState) {
	// TODO: Find out if rl.DrawImage would not be better for UI
	rl.DrawTextureEx(
		state.playerAvatar.Data,
		rl.Vector2{X: 10, Y: 490},
		0.0,  // rotation
		0.04, // scale
		rl.White,
	)

	r.drawTextOutlined("G demo", 10, 10, 20, rl.Black, rl.LightGray, 2)
	r.drawTextOutlined("Player hp:"+strconv.Itoa(state.playerHp), 160, 540, 20, rl.Black, rl.White, 2)
	r.drawHPBar(160, 570, 150, 12, state.playerHp, 100, rl.Black, rl.White)

}

func (r *Renderer) SetDebugMode(enabled bool) {
	r.debugMode = enabled
}

// Draw text with desired outline
// TODO: Move to UI utils
func (r *Renderer) drawTextOutlined(
	text string,
	x, y, fontSize int32,
	textColor, outlineColor rl.Color,
	thickness int32,
) {
	// Draw outline
	for ox := -thickness; ox <= thickness; ox++ {
		for oy := -thickness; oy <= thickness; oy++ {
			if ox != 0 && oy != 0 {
				rl.DrawText(text, x+ox, y+oy, fontSize, outlineColor)
			}
		}
	}

	rl.DrawText(text, x, y, fontSize, textColor)
}

// Draw HP bar of player
// TODO: Fix later, and move to proper place
func (r *Renderer) drawHPBar(
	x, y int32,
	width, height int32,
	current, max int,
	backColor, borderColor rl.Color,
) {

	// TODO: This should me somewhere else - drawing function is not to check limits of health
	if max <= 0 {
		return
	}

	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	pom := float32(current) / float32(max)

	var red, green uint8

	if pom >= 0.5 {
		pom = (1 - pom) * 2
		red = uint8(255 * pom)
		green = 255
	} else {
		red = 255
		pom = pom * 2
		green = uint8(255 * pom)
	}

	fillColor := rl.NewColor(red, green, 0, 255)
	percent := float32(current) / float32(max)
	fillWidth := int32(float32(width) * percent)

	rl.DrawRectangle(x-2, y-2, width+4, height+4, borderColor) // Border
	rl.DrawRectangle(x, y, width, height, backColor)           // Background
	rl.DrawRectangle(x, y, fillWidth, height, fillColor)       // Fill
}

func (r Renderer) renderLights(lights []entities.Light) {
	rl.DrawSphereEx(lights[0].Position, 0.2, 8, 8, rl.White)
	rl.DrawSphereEx(lights[1].Position, 0.2, 8, 8, rl.White)
	rl.DrawSphereEx(lights[2].Position, 0.2, 8, 8, rl.White)
}

// here you can put rendering basically everything, it can be turn off by debugMode flag
func (r *Renderer) renderDebug(state *GameState) {
	room := state.GetCurrentRoom()

	// rendering debug things from main.go
	if len(room.Objects) > 1 && room.Objects[1] != nil {
		r.renderCylinderSides(room.Objects[1], state.GetPlayerPosition())
	}

	// this neat thing used as celownik
	if point := state.GetRayCollisionPoint(); point != nil {
		r.renderCollisionPoint(*point)
	}
}

func (r *Renderer) renderCylinderSides(object *entities.Object, playerPos rl.Vector3) {
	cylinder, ok := object.Colliders[0].(*colliders.CylinderCollider)
	if !ok {
		return
	}

	// this things aroung debug cylinder
	pos1, pos2 := cylinder.GetSides(physics.GetVector2XZ(playerPos))

	drawPoint1 := rl.Vector3Add(physics.GetVector3FromXZ(pos1), cylinder.Position)
	drawPoint1.Y += 0.5

	drawPoint2 := rl.Vector3Add(physics.GetVector3FromXZ(pos2), cylinder.Position)
	drawPoint2.Y += 0.5

	// Create temporary debug objects for rendering
	debugCylinder1 := createDebugCylinder(drawPoint1)
	debugCylinder1.Draw()

	debugCylinder2 := createDebugCylinder(drawPoint2)
	debugCylinder2.Draw()
}

// util functions to render cube at ray collition point
func (r *Renderer) renderCollisionPoint(point rl.Vector3) {
	createDebugCube(rl.Vector3Add(point, rl.NewVector3(-0.05, -0.05, -0.05)), rl.Black).Draw()
}

// cube machine
func createDebugCube(position rl.Vector3, color rl.Color) entities.Object {
	return entities.Object{
		Model:     levels.NewModelFromCollider(colliders.NewCubeCollider(rl.Vector3{}, 0.1, 0.1, 0.1)),
		DrawPoint: position,
		Color:     color,
	}
}

// walec factory
func createDebugCylinder(position rl.Vector3) entities.Object {
	return entities.Object{
		Model:     levels.NewModelFromCollider(colliders.NewCylinderCollider(rl.Vector3{}, 0.1, 0.2)),
		DrawPoint: position,
		Color:     rl.Pink,
	}
}
