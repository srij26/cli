package shared

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	"code.cloudfoundry.org/cli/api/networkpolicy"
	"code.cloudfoundry.org/cli/api/networkpolicy/wrapper"
	"code.cloudfoundry.org/cli/api/uaa"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/translatableerror"
)

func NewClients(config command.Config) (*networkpolicy.Client, error) {
	ccClient := ccv2.NewClient(ccv2.Config{
		AppName:            config.BinaryName(),
		AppVersion:         config.BinaryVersion(),
		JobPollingTimeout:  config.OverallPollingTimeout(),
		JobPollingInterval: config.PollingInterval(),
		Wrappers:           nil,
	})

	if config.Target() == "" {
		return nil, translatableerror.NoAPISetError{
			BinaryName: config.BinaryName(),
		}
	}

	_, err := ccClient.TargetCF(ccv2.TargetSettings{
		URL:               config.Target(),
		SkipSSLValidation: config.SkipSSLValidation(),
		DialTimeout:       config.DialTimeout(),
	})
	if err != nil {
		return nil, err
	}

	uaaClient := uaa.NewClient(uaa.Config{
		AppName:           config.BinaryName(),
		AppVersion:        config.BinaryVersion(),
		ClientID:          config.UAAOAuthClient(),
		ClientSecret:      config.UAAOAuthClientSecret(),
		DialTimeout:       config.DialTimeout(),
		SkipSSLValidation: config.SkipSSLValidation(),
		URL:               ccClient.TokenEndpoint(),
	})

	networkpolicyClient := networkpolicy.NewClient(networkpolicy.Config{
		AppName:           config.BinaryName(),
		AppVersion:        config.BinaryVersion(),
		DialTimeout:       config.DialTimeout(),
		SkipSSLValidation: config.SkipSSLValidation(),
		URL:               "https://api.slouch.c2c.cf-app.com/networking/v0/external",
	})

	networkpolicyClient.WrapConnection(wrapper.NewUAAAuthentication(uaaClient, config))

	return networkpolicyClient, err
}
