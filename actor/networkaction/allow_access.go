package networkaction

import "code.cloudfoundry.org/cli/api/networkpolicy"

func (a *Actor) AddPolicy(sourceGUID, destinationGUID, protocol string, port int) error {
	return a.PolicyClient.AddPolicies([]networkpolicy.Policy{
		{
			Source: networkpolicy.Source{
				ID: sourceGUID,
			},
			Destination: networkpolicy.Destination{
				ID:       destinationGUID,
				Protocol: protocol,
				Port:     port,
			},
		},
	})
}
