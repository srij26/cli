package pushaction

import (
	"io"

	"code.cloudfoundry.org/cli/actor/v2action"
)

//go:generate counterfeiter . V2Actor

type V2Actor interface {
	BindRouteToApplication(routeGUID string, appGUID string) ([]string, error)
	CreateApplication(application v2action.Application) (v2action.Application, []string, error)
	CreateRoute(route v2action.Route, generatePort bool) (v2action.Route, []string, error)
	FindRouteBoundToSpaceWithSettings(route v2action.Route) (v2action.Route, []string, error)
	GatherDirectoryResources(sourceDir string) ([]v2action.Resource, error)
	GetApplicationByNameAndSpace(name string, spaceGUID string) (v2action.Application, []string, error)
	GetApplicationRoutes(applicationGUID string) ([]v2action.Route, []string, error)
	GetOrganizationDomains(orgGUID string) ([]v2action.Domain, []string, error)
	PollJob(job v2action.Job) ([]string, error)
	UpdateApplication(application v2action.Application) (v2action.Application, []string, error)
	UploadApplicationPackage(appGUID string, existingResources []v2action.Resource, newResources io.Reader, newResourcesLength int64) (v2action.Job, []string, error)
	ZipResources(sourceDir string, filesToInclude []v2action.Resource) (string, error)
}
