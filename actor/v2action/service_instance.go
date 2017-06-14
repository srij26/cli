package v2action

import (
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
)

// ServiceInstance represents an instance of a service.
type ServiceInstance ccv2.ServiceInstance

type ServiceInstanceNotFoundError struct {
	Name string
}

func (e ServiceInstanceNotFoundError) Error() string {
	return fmt.Sprintf("Service instance '%s' not found.", e.Name)
}

func (actor Actor) GetServiceInstanceByNameAndSpace(name string, spaceGUID string) (ServiceInstance, []string, error) {
	serviceInstances, warnings, err := actor.CloudControllerClient.GetSpaceServiceInstances(
		spaceGUID,
		true,
		[]ccv2.Query{
			ccv2.Query{
				Filter:   ccv2.NameFilter,
				Operator: ccv2.EqualOperator,
				Value:    name,
			},
		})

	if err != nil {
		return ServiceInstance{}, []string(warnings), err
	}

	if len(serviceInstances) == 0 {
		return ServiceInstance{}, []string(warnings), ServiceInstanceNotFoundError{
			Name: name,
		}
	}

	return ServiceInstance(serviceInstances[0]), []string(warnings), nil
}

func (actor Actor) GetServiceInstancesBySpace(spaceGUID string) ([]ServiceInstance, []string, error) {
	ccv2ServiceInstances, warnings, err := actor.CloudControllerClient.GetSpaceServiceInstances(
		spaceGUID, true, nil)

	if err != nil {
		return []ServiceInstance{}, []string(warnings), err
	}

	serviceInstances := make([]ServiceInstance, len(ccv2ServiceInstances))
	for i, ccv2ServiceInstance := range ccv2ServiceInstances {
		serviceInstances[i] = ServiceInstance(ccv2ServiceInstance)
	}

	return serviceInstances, []string(warnings), nil
}
