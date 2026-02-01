package server

import (
	"context"
	"fmt"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/entities"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/physics/colliders"
	"github.com/PawelZabc/ProjektZespolowy/internal/game/state"
)

type App struct {
	config    config.ServerConfig
	network   *Network
	gameState *GameState

	running bool
}

func NewApp(cfg config.ServerConfig) *App {
	return &App{
		config:  cfg,
		running: false,
	}
}

// Main game loop
func (a *App) updateLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(a.config.PhysicsTickRate))
	defer ticker.Stop()

	updateCount := int64(0)
	lastSend := 0
	sendRatio := float64(a.config.NetworkSendRate) / float64(a.config.PhysicsTickRate)

	for range ticker.C {
		if !a.running {
			break
		}

		updateCount++
		a.network.SetUpdateCount(updateCount)

		// Update game entities
		players := a.gameState.GetClientsAsPlayerSlice()
		a.updatePlayers(players, a.gameState.objects, a.gameState.enemy)
		a.updateEnemy(players, a.gameState.objects, a.gameState.enemy)

		a.network.RemoveDisconnectedClients(updateCount)

		// Sends updated game state to clients if it is time for that
		if (sendRatio*float64(updateCount))-float64(lastSend) >= 1 {
			lastSend++
			a.network.BroadcastGameState()
		}
	}
}

// Helper for updating Players. Handles physics (gravity, colliders) for players
func (a *App) updatePlayers(players Players, objects []colliders.Collider, enemy *entities.Enemy) {
	for _, player := range players {
		player.Velocity.Y -= config.Gravity

		player.Move()
		player.IsOnFloor = false

		for _, obj := range objects {
			player.PushbackFrom(obj)
		}

		player.PushbackFrom(enemy.Collider)
	}
}

// Helper for updating enemy.
// TODO: Change when there is more enemies
func (a *App) updateEnemy(players Players, objects []colliders.Collider, enemy *entities.Enemy) {
	enemy.Update(players, &objects)

	for _, obj := range objects {
		if obj != nil {
			enemy.Collider.PushbackFrom(obj)
		}
	}

	for _, player := range players {
		if enemy.Collider.PushbackFrom(player.Collider) != physics.DirNone {
			enemy.SetState(state.Attacking)
		}
	}
}

// Runs server application. Initialises app components. Starts receiving goroutine and game loop
func (a *App) Run() error {
	defer a.cleanup()

	if err := a.initComponents(); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	// Receiving goroutine - starting and ending (with context, to remove threads after the main thread is dead)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.network.StartReceiving(ctx)

	a.running = true
	a.updateLoop()

	return nil
}

// Initialises all of app components. Runs before game loop
func (a *App) initComponents() error {
	var err error

	a.gameState = NewGameState()
	a.network, err = NewNetwork(a.config.Port, a.gameState)
	if err != nil {
		return fmt.Errorf("network initialization failed: %w", err)
	}

	return nil
}

func (a *App) cleanup() {
	if a.network != nil {
		a.network.Close()
	}
}

func (a *App) Stop() {
	a.running = false
}
