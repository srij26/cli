package v2action

import "code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"

//go:generate counterfeiter . CloudControllerClient

// CloudControllerClient is a Cloud Controller V2 client.
type CloudControllerClient interface {
	AssociateSpaceWithSecurityGroup(securityGroupGUID string, spaceGUID string) ([]string, error)
	BindRouteToApplication(routeGUID string, appGUID string) (ccv2.Route, []string, error)
	CheckRoute(route ccv2.Route) (bool, []string, error)
	CreateApplication(app ccv2.Application) (ccv2.Application, []string, error)
	CreateRoute(route ccv2.Route, generatePort bool) (ccv2.Route, []string, error)
	CreateServiceBinding(appGUID string, serviceBindingGUID string, parameters map[string]interface{}) (ccv2.ServiceBinding, []string, error)
	CreateUser(uaaUserID string) (ccv2.User, []string, error)
	DeleteOrganization(orgGUID string) (ccv2.Job, []string, error)
	DeleteRoute(routeGUID string) ([]string, error)
	DeleteServiceBinding(serviceBindingGUID string) ([]string, error)
	GetApplication(guid string) (ccv2.Application, []string, error)
	GetApplicationInstancesByApplication(guid string) (map[int]ccv2.ApplicationInstance, []string, error)
	GetApplicationInstanceStatusesByApplication(guid string) (map[int]ccv2.ApplicationInstanceStatus, []string, error)
	GetApplicationRoutes(appGUID string, queries []ccv2.Query) ([]ccv2.Route, []string, error)
	GetApplications(queries []ccv2.Query) ([]ccv2.Application, []string, error)
	GetJob(jobGUID string) (ccv2.Job, []string, error)
	GetOrganization(guid string) (ccv2.Organization, []string, error)
	GetOrganizationPrivateDomains(orgGUID string, queries []ccv2.Query) ([]ccv2.Domain, []string, error)
	GetOrganizationQuota(guid string) (ccv2.OrganizationQuota, []string, error)
	GetOrganizations(queries []ccv2.Query) ([]ccv2.Organization, []string, error)
	GetPrivateDomain(domainGUID string) (ccv2.Domain, []string, error)
	GetRouteApplications(routeGUID string, queries []ccv2.Query) ([]ccv2.Application, []string, error)
	GetRoutes(queries []ccv2.Query) ([]ccv2.Route, []string, error)
	GetSecurityGroups(queries []ccv2.Query) ([]ccv2.SecurityGroup, []string, error)
	GetServiceBindings(queries []ccv2.Query) ([]ccv2.ServiceBinding, []string, error)
	GetServiceInstances(queries []ccv2.Query) ([]ccv2.ServiceInstance, []string, error)
	GetSharedDomain(domainGUID string) (ccv2.Domain, []string, error)
	GetSharedDomains() ([]ccv2.Domain, []string, error)
	GetSpaceQuota(guid string) (ccv2.SpaceQuota, []string, error)
	GetSpaceRoutes(spaceGUID string, queries []ccv2.Query) ([]ccv2.Route, []string, error)
	GetSpaceRunningSecurityGroupsBySpace(spaceGUID string) ([]ccv2.SecurityGroup, []string, error)
	GetSpaces(queries []ccv2.Query) ([]ccv2.Space, []string, error)
	GetSpacesBySecurityGroup(securityGroupGUID string) ([]ccv2.Space, []string, error)
	GetSpaceServiceInstances(spaceGUID string, includeUserProvidedServices bool, queries []ccv2.Query) ([]ccv2.ServiceInstance, []string, error)
	GetSpaceStagingSecurityGroupsBySpace(spaceGUID string) ([]ccv2.SecurityGroup, []string, error)
	GetStack(guid string) (ccv2.Stack, []string, error)
	PollJob(job ccv2.Job) ([]string, error)
	RemoveSpaceFromSecurityGroup(securityGroupGUID string, spaceGUID string) ([]string, error)
	TargetCF(settings ccv2.TargetSettings) ([]string, error)
	UpdateApplication(app ccv2.Application) (ccv2.Application, []string, error)
	UploadApplicationPackage(appGUID string, existingResources []ccv2.Resource, newResources ccv2.Reader, newResourcesLength int64) (ccv2.Job, []string, error)

	API() string
	APIVersion() string
	AuthorizationEndpoint() string
	DopplerEndpoint() string
	MinCLIVersion() string
	RoutingEndpoint() string
	TokenEndpoint() string
}
