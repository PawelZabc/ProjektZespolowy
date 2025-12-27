package server

import (
	"context"
	"fmt"
	"time"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
)

type App struct {
	config        config.ServerConfig
	network       *Network
	gameState     *GameState
	physicsEngine *Physics
	clientManager *ClientManager

	running bool
}

func NewApp(cfg config.ServerConfig) *App {
	return &App{
		config:  cfg,
		running: false,
	}
}

func (a *App) Run() error {

	defer a.cleanup()

	if err := a.initComponents(); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.network.StartReceiving(ctx)

	a.running = true
	a.updateLoop()

	return nil
}

func (a *App) initComponents() error {
	var err error

	a.gameState = NewGameState()
	a.clientManager = NewClientManager(a.gameState, a.config)
	a.network, err = NewNetwork(a.config.Port, a.clientManager, a.gameState)
	if err != nil {
		return fmt.Errorf("network initialization failed: %w", err)
	}
	a.physicsEngine = NewPhysics(a.gameState, a.clientManager)

	return nil
}

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

		a.physicsEngine.Update()
		a.clientManager.RemoveDisconnected(updateCount)
		if (sendRatio*float64(updateCount))-float64(lastSend) >= 1 {
			lastSend++
			a.network.BroadcastGameState()
		}
	}
}

func (a *App) cleanup() {
	if a.network != nil {
		a.network.Close()
	}
}

func (a *App) Stop() {
	a.running = false
}
