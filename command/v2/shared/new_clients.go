package shared

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	ccWrapper "code.cloudfoundry.org/cli/api/cloudcontroller/wrapper"
	"code.cloudfoundry.org/cli/api/uaa"
	uaaWrapper "code.cloudfoundry.org/cli/api/uaa/wrapper"
	"code.cloudfoundry.org/cli/command"
)

// NewClients creates a new V2 Cloud Controller client and UAA client using the
// passed in config.
func NewClients(config command.Config, ui command.UI) (*ccv2.Client, *uaa.Client, error) {
	if config.Target() == "" {
		return nil, nil, command.NoAPISetError{
			BinaryName: config.BinaryName(),
		}
	}

	ccClient := ccv2.NewClient(ccv2.Config{
		AppName:            config.BinaryName(),
		AppVersion:         config.BinaryVersion(),
		JobPollingTimeout:  config.OverallPollingTimeout(),
		JobPollingInterval: config.PollingInterval(),
		DialTimeout:        config.DialTimeout(),
		SkipSSLValidation:  config.SkipSSLValidation(),
	})

	ccClient.WrapConnection(ccv2.NewErrorWrapper()) //Pretty Sneaky, Sis..

	verbose, location := config.Verbose()
	if verbose {
		ccClient.WrapConnection(ccWrapper.NewRequestLogger(ui.RequestLoggerTerminalDisplay()))
	}
	if location != nil {
		ccClient.WrapConnection(ccWrapper.NewRequestLogger(ui.RequestLoggerFileWriter(location)))
	}

	ccUAAWrapper := ccWrapper.NewUAAAuthentication(nil, config)

	ccClient.WrapConnection(ccUAAWrapper)
	ccClient.WrapConnection(ccWrapper.NewRetryRequest(2))

	_, err := ccClient.TargetCF(config.Target())
	if err != nil {
		return nil, nil, err
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

	uaaClient.WrapConnection(uaa.NewErrorWrapper())

	if verbose {
		uaaClient.WrapConnection(uaaWrapper.NewRequestLogger(ui.RequestLoggerTerminalDisplay()))
	}

	if location != nil {
		uaaClient.WrapConnection(uaaWrapper.NewRequestLogger(ui.RequestLoggerFileWriter(location)))
	}

	uaaUAAWrapper := uaaWrapper.NewUAAAuthentication(uaaClient, config)

	uaaClient.WrapConnection(uaaUAAWrapper)
	uaaClient.WrapConnection(uaaWrapper.NewRetryRequest(2))

	ccUAAWrapper.SetClient(uaaClient)

	return ccClient, uaaClient, err
}
