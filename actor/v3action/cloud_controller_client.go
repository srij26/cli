package v3action

import (
	"net/url"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

//go:generate counterfeiter . CloudControllerClient

// CloudControllerClient is the interface to the cloud controller V3 API.
type CloudControllerClient interface {
	AssignSpaceToIsolationSegment(spaceGUID string, isolationSegmentGUID string) (ccv3.Relationship, []string, error)
	CloudControllerAPIVersion() string
	CreateApplication(app ccv3.Application) (ccv3.Application, []string, error)
	CreateApplicationTask(appGUID string, task ccv3.Task) (ccv3.Task, []string, error)
	CreateBuild(build ccv3.Build) (ccv3.Build, []string, error)
	CreateIsolationSegment(isolationSegment ccv3.IsolationSegment) (ccv3.IsolationSegment, []string, error)
	CreatePackage(pkg ccv3.Package) (ccv3.Package, []string, error)
	DeleteIsolationSegment(guid string) ([]string, error)
	EntitleIsolationSegmentToOrganizations(isoGUID string, orgGUIDs []string) (ccv3.RelationshipList, []string, error)
	GetApplications(query url.Values) ([]ccv3.Application, []string, error)
	GetApplicationCurrentDroplet(appGUID string) (ccv3.Droplet, []string, error)
	GetApplicationProcesses(appGUID string) ([]ccv3.Process, []string, error)
	GetProcessInstances(processGUID string) ([]ccv3.Instance, []string, error)
	GetApplicationTasks(appGUID string, query url.Values) ([]ccv3.Task, []string, error)
	GetBuild(guid string) (ccv3.Build, []string, error)
	GetIsolationSegment(guid string) (ccv3.IsolationSegment, []string, error)
	GetIsolationSegmentOrganizationsByIsolationSegment(isolationSegmentGUID string) ([]ccv3.Organization, []string, error)
	GetIsolationSegments(query url.Values) ([]ccv3.IsolationSegment, []string, error)
	GetOrganizationDefaultIsolationSegment(orgGUID string) (ccv3.Relationship, []string, error)
	GetOrganizations(query url.Values) ([]ccv3.Organization, []string, error)
	GetPackage(guid string) (ccv3.Package, []string, error)
	GetSpaceIsolationSegment(spaceGUID string) (ccv3.Relationship, []string, error)
	RevokeIsolationSegmentFromOrganization(isolationSegmentGUID string, organizationGUID string) ([]string, error)
	SetApplicationDroplet(appGUID string, dropletGUID string) (ccv3.Relationship, []string, error)
	StartApplication(appGUID string) (ccv3.Application, []string, error)
	UpdateTask(taskGUID string) (ccv3.Task, []string, error)
	UploadPackage(pkg ccv3.Package, zipFilepath string) (ccv3.Package, []string, error)
}
