package config

// struct to keep client config
type ClientConfig struct {
	ServerIP   string
	ServerPort int

	WindowWidth  int
	WindowHeight int
	WindowTitle  string
	TargetFPS    int

	DebugMode bool
}

// to test using localhost - it will be shadowed in prod
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		ServerIP:     "127.0.0.1",
		ServerPort:   DefaultPort,
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "MGT - Maybe Game Title",
		TargetFPS:    60,
		DebugMode:    true,
	}
}
