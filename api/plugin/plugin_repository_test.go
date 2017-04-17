package plugin_test

import (
	"net/http"

	. "code.cloudfoundry.org/cli/api/plugin"
	"code.cloudfoundry.org/cli/api/plugin/pluginerror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("PluginRepository", func() {
	var client *Client

	BeforeEach(func() {
		client = NewTestClient()
	})

	Describe("GetPluginRepository", func() {
		Context("when the url is a CF CLI plugin repo", func() {
			BeforeEach(func() {
				response := `{
					"plugins": [
						{
							"name": "plugin-1",
							"description": "useful plugin for useful things",
							"version": "1.0.0"
						},
						{
							"name": "plugin-2",
							"description": "amazing plugin",
							"version": "1.0.0"
						}
					]
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/list"),
						RespondWith(http.StatusOK, response),
					),
				)
			})

			It("returns the plugin repository", func() {
				pluginRepository, err := client.GetPluginRepository(testPluginRepoURL())
				Expect(err).ToNot(HaveOccurred())
				Expect(pluginRepository).To(Equal(PluginRepository{
					Plugins: []Plugin{
						{
							Name:        "plugin-1",
							Description: "useful plugin for useful things",
							Version:     "1.0.0",
						},
						{
							Name:        "plugin-2",
							Description: "amazing plugin",
							Version:     "1.0.0",
						},
					},
				}))
			})
		})

		Context("when the http client returns an error", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/list"),
						RespondWith(http.StatusNotFound, nil),
					),
				)
			})

			It("returns the error", func() {
				_, err := client.GetPluginRepository(testPluginRepoURL())
				Expect(err).To(MatchError(pluginerror.NotFoundError{}))
			})
		})
	})
})
