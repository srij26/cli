package networkpolicy

import (
	"code.cloudfoundry.org/cli/actor/networkaction"
	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/networkpolicy/shared"
	v2Shared "code.cloudfoundry.org/cli/command/v2/shared"
)

//go:generate counterfeiter . AllowAccessActor

type AllowAccessActor interface {
	AddPolicy(sourceApp, destinationApp, protocol string, port int) error
}

//go:generate counterfeiter . AllowAccessCCActor

type AllowAccessCCActor interface {
	GetApplicationByNameAndSpace(name string, spaceGUID string) (v2action.Application, v2action.Warnings, error)
}

type AllowAccessCommand struct {
	RequiredArgs flag.AllowAccessArgs `positional-args:"yes"`
	Protocol     string               `long:"protocol" description:"Protocol to connect apps with. (required)"`
	Port         int                  `long:"port" description:"Port to connect to destination app with. (required)"`
	usage        interface{}          `usage:"CF_NAME allow-access SOURCE_APP DESTINATION_APP --protocol <tcp|udp> --port <1-65535>"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	CCActor     AllowAccessCCActor
	Actor       AllowAccessActor
}

func (cmd *AllowAccessCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor()

	ccClient, uaaClient, err := v2Shared.NewClients(config, ui, true)
	if err != nil {
		return err
	}
	cmd.CCActor = v2action.NewActor(ccClient, uaaClient)

	networkpolicyClient, err := shared.NewClients(config)
	if err != nil {
		return err
	}
	cmd.Actor = networkaction.NewActor(networkpolicyClient)

	return nil
}

func (cmd AllowAccessCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(cmd.Config, true, true)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	sourceApp, _, err := cmd.CCActor.GetApplicationByNameAndSpace(cmd.RequiredArgs.SourceApp, cmd.Config.TargetedSpace().GUID)
	if err != nil {
		return err
	}

	destinationApp, _, err := cmd.CCActor.GetApplicationByNameAndSpace(cmd.RequiredArgs.DestinationApp, cmd.Config.TargetedSpace().GUID)
	if err != nil {
		return err
	}

	cmd.UI.DisplayText("Allowing traffic from {{.SOURCE_APP}} to {{.DESTINATION_APP}} as {{.CURRENT_USER}}...", map[string]interface{}{
		"SOURCE_APP":      cmd.RequiredArgs.SourceApp,
		"DESTINATION_APP": cmd.RequiredArgs.DestinationApp,
		"CURRENT_USER":    user.Name,
	})

	err = cmd.Actor.AddPolicy(sourceApp.GUID, destinationApp.GUID, cmd.Protocol, cmd.Port)
	if err != nil {
		return err
	}

	cmd.UI.DisplayOK()

	return nil
}
