package configv3

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// PluginsConfig represents the plugin configuration
type PluginsConfig struct {
	Plugins map[string]Plugin `json:"Plugins"`
}

// Plugin represents the plugin as a whole, not be confused with PluginCommand
type Plugin struct {
	// repo, filepath, url
	Name     string
	Location string          `json:"Location"`
	Version  PluginVersion   `json:"Version"`
	Commands []PluginCommand `json:"Commands"`
}

// PluginVersion is the plugin version information
type PluginVersion struct {
	Major int `json:"Major"`
	Minor int `json:"Minor"`
	Build int `json:"Build"`
}

// PluginCommand represents an individual command inside a plugin
type PluginCommand struct {
	Name         string             `json:"Name"`
	Alias        string             `json:"Alias"`
	HelpText     string             `json:"HelpText"`
	UsageDetails PluginUsageDetails `json:"UsageDetails"`
}

// PluginUsageDetails contains the usage metadata provided by the plugin
type PluginUsageDetails struct {
	Usage   string            `json:"Usage"`
	Options map[string]string `json:"Options"`
}

// Plugins returns back the plugin configuration read from the plugin home
func (config *Config) Plugins() []Plugin {
	plugins := []Plugin{}
	for _, plugin := range config.pluginConfig.Plugins {
		plugins = append(plugins, plugin)
	}
	sort.Slice(plugins, func(i, j int) bool {
		return strings.ToLower(plugins[i].Name) < strings.ToLower(plugins[j].Name)
	})
	return plugins
}

// CalculateSHA1 returns the sha1 value of the plugin executable. If an error
// is encountered calculating sha1, N/A is returned
func (p Plugin) CalculateSHA1() string {
	fileSHA := ""
	contents, err := ioutil.ReadFile(p.Location)
	if err != nil {
		fileSHA = "N/A"
	} else {
		fileSHA = fmt.Sprintf("%x", sha1.Sum(contents))
	}
	return fileSHA
}

// PluginCommands returns the plugin's commands sorted by command name.
func (p Plugin) PluginCommands() []PluginCommand {
	sort.Slice(p.Commands, func(i, j int) bool {
		return strings.ToLower(p.Commands[i].Name) < strings.ToLower(p.Commands[j].Name)
	})
	return p.Commands
}

// String returns the plugin's version in the format x.y.z.
func (v PluginVersion) String() string {
	if v.Major == 0 && v.Minor == 0 && v.Build == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}

// CommandName returns the name of the plugin. The name is concatenated with
// alias if alias is specified.
func (c PluginCommand) CommandName() string {
	if c.Name != "" && c.Alias != "" {
		return fmt.Sprintf("%s, %s", c.Name, c.Alias)
	}
	return c.Name
}

// PluginHome returns the plugin configuration directory based off:
//   1. The $CF_PLUGIN_HOME environment variable if set
//   2. Defaults to the home diretory (outlined in LoadConfig)/.cf/plugins
func (config *Config) PluginHome() string {
	if config.ENV.CFPluginHome != "" {
		return filepath.Join(config.ENV.CFPluginHome, ".cf", "plugins")
	}

	return filepath.Join(homeDirectory(), ".cf", "plugins")
}
