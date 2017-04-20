package pluginaction

import "code.cloudfoundry.org/cli/util/configv3"

//go:generate counterfeiter . Config

// Config a way of getting basic CF configuration
type Config interface {
	PluginRepositories() []configv3.PluginRepository
	Plugins() []configv3.Plugin
}
