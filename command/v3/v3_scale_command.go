package v3

import (
	"strconv"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v3/shared"
	"github.com/cloudfoundry/bytefmt"
)

//go:generate counterfeiter . V3ScaleActor

type V3ScaleActor interface {
	GetProcessByApplicationNameAndSpace(appName string, spaceGUID string) (v3action.Process, v3action.Warnings, error)
	ScaleProcessByApplicationNameAndSpace(appName string, spaceGUID string, process ccv3.Process) (v3action.Warnings, error)
}

type V3ScaleCommand struct {
	RequiredArgs    flag.AppName   `positional-args:"yes"`
	Instances       int            `short:"i" description:"Number of instances"`
	DiskLimit       flag.Megabytes `short:"k" description:"Disk limit (e.g. 256M, 1024M, 1G)"`
	MemoryLimit     flag.Megabytes `short:"m" description:"Memory limit (e.g. 256M, 1024M, 1G)"`
	usage           interface{}    `usage:"CF_NAME v3-scale APP_NAME [-i INSTANCES] [-k DISK] [-m MEMORY]"`
	relatedCommands interface{}    `related_commands:"v3-push"`

	UI          command.UI
	Config      command.Config
	Actor       V3ScaleActor
	SharedActor command.SharedActor
}

func (cmd *V3ScaleCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor()

	ccClient, _, err := shared.NewClients(config, ui, true)
	if err != nil {
		return err
	}
	cmd.Actor = v3action.NewActor(ccClient, config)

	return nil
}

func (cmd V3ScaleCommand) Execute(args []string) error {
	cmd.UI.DisplayText(command.ExperimentalWarning)
	cmd.UI.DisplayNewline()

	err := cmd.SharedActor.CheckTarget(cmd.Config, true, true)
	if err != nil {
		return shared.HandleError(err)
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	var actorErr error
	if cmd.Instances == 0 && cmd.DiskLimit.Size == 0 && cmd.MemoryLimit.Size == 0 {
		actorErr = cmd.displayProcessInformation(user.Name)
	} else {
		actorErr = cmd.scaleProcess(user.Name)
	}
	if actorErr != nil {
		return actorErr
	}

	cmd.UI.DisplayOK()

	return nil
}

func (cmd V3ScaleCommand) displayProcessInformation(username string) error {
	cmd.UI.DisplayTextWithFlavor("Showing current scale of app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...", map[string]interface{}{
		"AppName":   cmd.RequiredArgs.AppName,
		"OrgName":   cmd.Config.TargetedOrganization().Name,
		"SpaceName": cmd.Config.TargetedSpace().Name,
		"Username":  username,
	})

	process, warnings, err := cmd.Actor.GetProcessByApplicationNameAndSpace(cmd.RequiredArgs.AppName, cmd.Config.TargetedSpace().GUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayNewline()
	cmd.UI.DisplayKeyValueTable("", [][]string{
		{cmd.UI.TranslateText("memory:"), bytefmt.ByteSize(uint64(process.MemoryInMB) * bytefmt.MEGABYTE)},
		{cmd.UI.TranslateText("disk:"), bytefmt.ByteSize(uint64(process.DiskInMB) * bytefmt.MEGABYTE)},
		{cmd.UI.TranslateText("instances:"), strconv.Itoa(len(process.Instances))},
	}, 3)

	return nil
}

func (cmd V3ScaleCommand) scaleProcess(username string) error {
	cmd.UI.DisplayTextWithFlavor("Scaling app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...", map[string]interface{}{
		"AppName":   cmd.RequiredArgs.AppName,
		"OrgName":   cmd.Config.TargetedOrganization().Name,
		"SpaceName": cmd.Config.TargetedSpace().Name,
		"Username":  username,
	})

	ccv3Process := ccv3.Process{
		Type:       "web",
		Instances:  cmd.Instances,
		MemoryInMB: int(cmd.MemoryLimit.Size),
		DiskInMB:   int(cmd.DiskLimit.Size),
	}
	warnings, err := cmd.Actor.ScaleProcessByApplicationNameAndSpace(cmd.RequiredArgs.AppName, cmd.Config.TargetedSpace().GUID, ccv3Process)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	return nil
}
