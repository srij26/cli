package pluginaction_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPluginaction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pluginaction Suite")
}
