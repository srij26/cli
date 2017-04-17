package plugin

// PluginRepository represents a plugin repository
type PluginRepository struct {
	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

func (client *Client) GetPluginRepository(repositoryURL string) (PluginRepository, error) {
	request, err := client.newHTTPGetRequest(repositoryURL)
	if err != nil {
		return PluginRepository{}, err
	}

	var pluginRepository PluginRepository
	response := Response{
		Result: &pluginRepository,
	}
	err = client.connection.Make(request, &response)
	if err != nil {
		return PluginRepository{}, err
	}

	return pluginRepository, nil
}
