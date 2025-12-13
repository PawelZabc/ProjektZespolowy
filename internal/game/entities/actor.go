package entities

import (
	"github.com/PawelZabc/ProjektZespolowy/assets"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Actor struct {
	levels.Object
	Rotation float32
	// Animation
}

func (a *Actor) SetPosition(pos rl.Vector3) {
	a.Object.Colliders[0].SetPosition(pos)
}

func DrawActors(actors []*Actor) {
	for _, actor := range actors {
		actor.Draw()
	}
}

func DrawActorsMap[T comparable](actors map[T]*Actor) {
	for _, actor := range actors {
		actor.Draw()
	}
}

func (a *Actor) Draw() {
	rl.DrawModelEx(a.Object.Model, rl.Vector3Add(a.Object.DrawPoint, a.Colliders[0].GetPosition()), rl.NewVector3(0, 1, 0), a.Rotation, rl.Vector3One(), a.Object.Color)
}

func NewActor(collider colliders.Collider, drawPoint rl.Vector3, rotation float32, modelName string) *Actor {
	// fmt.Println("Pre Load Actor Model !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	model, _ := assets.GlobalManager.LoadModel(modelName)
	// fmt.Println("Load Actor Model Error", err)
	object := levels.Object{Model: model.Data, DrawPoint: drawPoint, Color: rl.White, Colliders: []colliders.Collider{collider}}
	actor := Actor{Object: object, Rotation: rotation}
	return &actor
}
