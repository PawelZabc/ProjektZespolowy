package config

// structo to hold server configuraiton
type ServerConfig struct {
	Port int

	PhysicsTickRate int // game loop updates
	NetworkSendRate int // how many tps info is sent to players
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:            DefaultPort,
		PhysicsTickRate: 60,
		NetworkSendRate: 30,
	}
}
