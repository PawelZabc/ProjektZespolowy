package client

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Renderer struct {
	// Go there if you want to draw things in game

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

	if r.debugMode {
		r.renderDebug(state)
	}
}

// randering 2d elements
func (r *Renderer) RenderUI() {
	rl.DrawText("Press 'R' to unlock mouse", 10, 10, 20, rl.Black)
	// here put other Ui thingis

}

func (r *Renderer) SetDebugMode(enabled bool) {
	r.debugMode = enabled
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
	createDebugCube(rl.Vector3Add(point, rl.NewVector3(-0.05, -0.05, -0.05))).Draw()
}

// cube machine
func createDebugCube(position rl.Vector3) entities.Object {
	return entities.Object{
		Model:     levels.NewModelFromCollider(colliders.NewCubeCollider(rl.Vector3{}, 0.1, 0.1, 0.1)),
		DrawPoint: position,
		Color:     rl.Pink,
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
