package pushaction

import (
	"fmt"

	"code.cloudfoundry.org/cli/actor/v2action"
	log "github.com/sirupsen/logrus"
)

// NoDomainsFoundError is returned when there are no private or shared domains
// accessible to an organization.
type NoDomainsFoundError struct {
	OrganizationGUID string
}

func (e NoDomainsFoundError) Error() string {
	return fmt.Sprintf("No private or shared domains found for organization (GUID: %s)", e.OrganizationGUID)
}

// DefaultDomain looks up the shared and then private domains and returns back
// the first one in the list as the default.
func (actor Actor) DefaultDomain(orgGUID string) (v2action.Domain, []string, error) {
	log.Infoln("getting org domains for org GUID:", orgGUID)
	domains, warnings, err := actor.V2Actor.GetOrganizationDomains(orgGUID)

	if err != nil {
		log.Errorln("searching for domains in org:", err)
		return v2action.Domain{}, []string(warnings), err
	}

	if len(domains) == 0 {
		log.Error("no domains found")
		return v2action.Domain{}, []string(warnings), NoDomainsFoundError{OrganizationGUID: orgGUID}
	}

	log.Debugf("selecting first domain as default domain: %#v", domains)
	return domains[0], []string(warnings), nil
}
