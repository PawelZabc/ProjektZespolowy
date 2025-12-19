package server

import (
	"fmt"

	"github.com/PawelZabc/ProjektZespolowy/internal/config"
)

type App struct {
	config  config.ServerConfig
	network *Network
	running bool
}

func NewApp(cfg config.ServerConfig) *App {
	return &App{
		config:  cfg,
		running: false,
	}
}

func (a *App) Run() error {

	defer a.cleanup() // end with cleanup

	// init components like input, network, state, etc..
	if err := a.initComponents(); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	a.running = true

	return nil
}

func (a *App) initComponents() error {
	var err error

	a.network, err = NewNetwork(a.config.Port)
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
