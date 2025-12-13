package entities

import (
	"github.com/PawelZabc/ProjektZespolowy/client/assets"
	"github.com/PawelZabc/ProjektZespolowy/client/game"
	types "github.com/PawelZabc/ProjektZespolowy/shared/_types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Actor struct {
	game.Object
	Rotation       float32
	Animation      uint8
	AnimationFrame uint8
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

func (a *Actor) SetAnimation(anim uint8) {
	if anim == a.Animation {
		a.AnimationFrame++
	} else {
		a.Animation = anim
		a.AnimationFrame = 0
	}
}

func (a *Actor) UpdateAnimation() { //change this when the actual animation gets added
	if a.Animation == uint8(types.Attacking) {
		notRed := 255 - (8 * a.AnimationFrame)
		if a.AnimationFrame > 50 {
			notRed = 0
		}
		a.Color = rl.NewColor(255, notRed, notRed, 255)
	} else {
		a.Color = rl.White
	}
}

func (a *Actor) Draw() {
	a.UpdateAnimation()
	rl.DrawModelEx(a.Object.Model, rl.Vector3Add(a.Object.DrawPoint, a.Colliders[0].GetPosition()), rl.NewVector3(0, 1, 0), a.Rotation, rl.Vector3One(), a.Object.Color)
}

func NewActor(collider types.Collider, drawPoint rl.Vector3, rotation float32, modelName string) *Actor {
	// fmt.Println("Pre Load Actor Model !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	model, _ := assets.GlobalManager.LoadModel(modelName)
	// fmt.Println("Load Actor Model Error", err)
	object := game.Object{Model: model.Data, DrawPoint: drawPoint, Color: rl.White, Colliders: []types.Collider{collider}}
	actor := Actor{Object: object, Rotation: rotation}
	return &actor
}
