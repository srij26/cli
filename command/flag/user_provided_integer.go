package flag

import (
	"strings"

	"github.com/cloudfoundry/bytefmt"
	flags "github.com/jessevdk/go-flags"
)

type UserProvidedInteger struct {
}

func (m *Megabytes) UnmarshalFlag(val string) error {
	return nil
}
