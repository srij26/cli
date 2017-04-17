package pluginaction

import (
	"fmt"
	"os"
)

// PluginNotFoundError is an error returned when a plugin is not found.
type PluginNotFoundError struct {
	Name string
}

// Error outputs a plugin not found error message.
func (e PluginNotFoundError) Error() string {
	return fmt.Sprintf("Plugin name %s does not exist", e.Name)
}

//go:generate counterfeiter . PluginUninstaller

type PluginUninstaller interface {
	Uninstall(pluginPath string) error
}

func (actor Actor) UninstallPlugin(uninstaller PluginUninstaller, name string) error {
	plugin, exist := actor.config.Plugins()[name]
	if !exist {
		return PluginNotFoundError{Name: name}
	}

	err := uninstaller.Uninstall(plugin.Location)
	if err != nil {
		return err
	}

	// sleep for 500 ms???
	// the plugin could potentially make API calls. we should think of a better strategy for how this plays with removing the plugin binary in the next step

	err = os.Remove(plugin.Location)
	if err != nil {
		return err
	}

	actor.config.RemovePlugin(name)
	return nil
}
