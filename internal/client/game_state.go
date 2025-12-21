package client

import (
	"fmt"
	"log"
	"sync"

	"github.com/PawelZabc/ProjektZespolowy/assets"
	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/levels"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/protocol"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState struct {
	player  *entities.Object           // local player
	players map[uint16]*entities.Actor // other players
	enemy   *entities.Actor            // for now only one

	rooms       []levels.ClientRoom
	currentRoom int
	shader      rl.Shader
	lights      []entities.Light

	rayCollisionPoint *rl.Vector3
	mu                sync.RWMutex
}

func NewGameState() *GameState {
	playerCollider := colliders.NewCylinderCollider(config.PlayerSpawnpoint, config.PlayerRadius, config.PlayerHeight)
	player := &entities.Object{
		Colliders: []colliders.Collider{playerCollider},
		Model:     levels.NewModelFromCollider(playerCollider),
	}

	enemy := entities.NewActor(
		colliders.NewCylinderCollider(rl.NewVector3(15, 0, 15), 1, 2),
		rl.Vector3{},
		-45,
		assets.ModelGhost,
	)

	// TODO: Move this shady code to somewhere else
	shaderPointer, err := assets.GlobalManager.LoadShader(assets.ShaderLightingV2VS, assets.ShaderLightingV2FS)

	if err != nil {
		log.Println(err)
	}

	shader := shaderPointer.Data

	// Set view position uniform
	*shader.Locs = rl.GetShaderLocation(shader, "viewPos")

	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	ambient := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(shader, ambientLoc, ambient, rl.ShaderUniformVec4)

	rooms := levels.LoadRooms(shader)

	lights := createLights(shader)

	// loading player model into memory - TODO: FIX
	// due to bug with loading model in another thread (OpenGL context is limited to one)
	playerModel, _ := assets.GlobalManager.LoadModel(assets.ModelPlayer)
	fmt.Println("Player model loaded", playerModel)

	return &GameState{
		player:      player,
		players:     make(map[uint16]*entities.Actor),
		enemy:       enemy,
		rooms:       rooms,
		shader:      shader,
		lights:      lights,
		currentRoom: 0,
	}
}

// TODO: Implement that method if it would be usefull
// Updates that happen only on client
func (gs *GameState) UpdateFromClient() {
	log.Fatal("UpdateFromClient not implemented")
}

// Gather data from server and apply it on the client (executes in goroutine)
func (gs *GameState) UpdateFromServer(data protocol.ServerData) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.player.Colliders[0].SetPosition(data.Position)

	gs.enemy.SetPosition(data.Enemy.Position)
	gs.enemy.Rotation = -data.Enemy.Rotation // why minus tho?

	updatedPlayers := make(map[uint16]bool)

	for _, playerData := range data.Players {
		if actor, exists := gs.players[playerData.Id]; exists {
			// update existing
			actor.Object.Colliders[0].SetPosition(playerData.Position)
			actor.Rotation = (playerData.Rotation * rl.Rad2deg)
		} else {
			// new player
			gs.createPlayer(playerData.Id, playerData.Position, playerData.Rotation)
		}
		updatedPlayers[playerData.Id] = true
	}

	// if player wasnt updated,its no longer at the server so delete it
	for id := range gs.players {
		if !updatedPlayers[id] {
			delete(gs.players, id)
		}
	}
}

func (gs *GameState) GetShader() rl.Shader {
	return gs.shader
}

func createLights(shader rl.Shader) []entities.Light {
	var lights []entities.Light

	light1 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(0, 2.9, 0), // light position
		rl.NewVector3(0, 0, 0),   // target (unused for point light)
		rl.White,                 // light color
		shader,
	)
	light1.UpdateValues()

	light2 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(5, 2.9, -5), // light position
		rl.NewVector3(0, 0, 0),    // target (unused for point light)
		rl.White,                  // light color
		shader,
	)
	light2.UpdateValues()

	light3 := entities.NewLight(
		entities.LightTypePoint,
		rl.NewVector3(5, 2.9, 5), // light position
		rl.NewVector3(0, 0, 0),   // target (unused for point light)
		rl.White,                 // light color
		shader,
	)
	light3.Enabled = 1 // ON is default
	light3.UpdateValues()

	lights = append(lights, light1, light2, light3)

	return lights
}

// Creates a new remote player
func (gs *GameState) createPlayer(id uint16, position rl.Vector3, rotation float32) {
	// maybe some logging could be usefull
	cylinder := colliders.NewCylinderCollider(position, config.PlayerRadius, config.PlayerHeight)
	gs.players[id] = entities.NewActor(
		cylinder,
		rl.Vector3{},
		(rotation*rl.Rad2deg)+90,
		assets.ModelPlayer,
	)
}

// Local player position
func (gs *GameState) GetPlayerPosition() rl.Vector3 {
	return gs.player.Colliders[0].GetPosition()
}

func (gs *GameState) GetCurrentRoom() *levels.ClientRoom {
	return &gs.rooms[gs.currentRoom]
}

// Getter for other players
func (gs *GameState) GetPlayers() map[uint16]*entities.Actor {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.players
}

// TODO: refactor it to get all enemies in the future
func (gs *GameState) GetEnemy() *entities.Actor {
	return gs.enemy
}

// Updates ray collision detection
func (gs *GameState) UpdateRayCollision(ray colliders.Ray) {
	var minLength float32 = 0
	gs.rayCollisionPoint = nil

	for _, object := range gs.rooms[gs.currentRoom].Objects {
		if object != nil {
			for _, collider := range object.Colliders {
				point, length := ray.GetCollisionPoint(collider)
				if point != nil && (minLength == 0 || length < minLength) {
					minLength = length
					gs.rayCollisionPoint = point
				}
			}
		}
	}
}

// GetRayCollisionPoint returns the current ray collision point
func (gs *GameState) GetRayCollisionPoint() *rl.Vector3 {
	return gs.rayCollisionPoint
}
