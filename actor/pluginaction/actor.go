package pluginaction

// Actor handles all shared actions
type Actor struct {
	config Config
	client PluginClient
}

// NewActor returns an Actor with default settings
func NewActor(config Config, client PluginClient) Actor {
	return Actor{config: config, client: client}
}
