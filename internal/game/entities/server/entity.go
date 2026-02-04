package server

import (
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
)

type SEntity interface {
	GetCollider() *colliders.Collider

	// Collider colliders.Collider
}

type EffectObject struct {
	Collider colliders.Collider
	Effect   *func(playerIp string)
	Active   bool
}

type CollisionObject struct {
	Collider *colliders.Collider
}
