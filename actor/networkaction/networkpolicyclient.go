package networkaction

import "code.cloudfoundry.org/cli/api/networkpolicy"

type NetworkPolicyClient interface {
	// GetPolicies(token string) ([]models.Policy, error)
	// GetPoliciesByID(token string, ids ...string) ([]models.Policy, error)
	// DeletePolicies(token string, policies []models.Policy) error
	AddPolicies(policies []networkpolicy.Policy) error
}
