package ccv3_test

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	. "code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Process", func() {
	var client *Client

	BeforeEach(func() {
		client = NewTestClient()
	})

	Describe("GetApplicationProcesses", func() {
		Context("when the application exists", func() {
			BeforeEach(func() {
				response1 := fmt.Sprintf(`
					{
						"pagination": {
							"next": {
								"href": "%s/v3/apps/some-app-guid/processes?page=2"
							}
						},
						"resources": [
							{
								"guid": "process-1-guid",
								"type": "web",
								"memory_in_mb": 32,
								"health_check": {
                  "type": "port",
                  "data": {
                    "timeout": null,
                    "endpoint": null
                  }
                }
							},
							{
								"guid": "process-2-guid",
								"type": "worker",
								"memory_in_mb": 64,
								"health_check": {
                  "type": "http",
                  "data": {
                    "timeout": 60,
                    "endpoint": "/health"
                  }
                }
							}
						]
					}`, server.URL())
				response2 := `
					{
						"pagination": {
							"next": null
						},
						"resources": [
							{
								"guid": "process-3-guid",
								"type": "console",
								"memory_in_mb": 128,
								"health_check": {
                  "type": "process",
                  "data": {
                    "timeout": 90,
                    "endpoint": null
                  }
                }
							}
						]
					}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes"),
						RespondWith(http.StatusOK, response1, http.Header{"X-Cf-Warnings": {"warning-1"}}),
					),
				)
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes", "page=2"),
						RespondWith(http.StatusOK, response2, http.Header{"X-Cf-Warnings": {"warning-2"}}),
					),
				)
			})

			It("returns a list of processes associated with the application and all warnings", func() {
				processes, warnings, err := client.GetApplicationProcesses("some-app-guid")
				Expect(err).ToNot(HaveOccurred())

				Expect(processes).To(ConsistOf(
					Process{
						GUID:        "process-1-guid",
						Type:        "web",
						MemoryInMB:  32,
						HealthCheck: ProcessHealthCheck{Type: "port"},
					},
					Process{
						GUID:       "process-2-guid",
						Type:       "worker",
						MemoryInMB: 64,
						HealthCheck: ProcessHealthCheck{
							Type: "http",
							Data: ProcessHealthCheckData{Endpoint: "/health"},
						},
					},
					Process{
						GUID:        "process-3-guid",
						Type:        "console",
						MemoryInMB:  128,
						HealthCheck: ProcessHealthCheck{Type: "process"},
					},
				))
				Expect(warnings).To(ConsistOf("warning-1", "warning-2"))
			})
		})

		Context("when cloud controller returns an error", func() {
			BeforeEach(func() {
				response := `{
					"errors": [
						{
							"code": 10010,
							"detail": "App not found",
							"title": "CF-ResourceNotFound"
						}
					]
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes"),
						RespondWith(http.StatusNotFound, response),
					),
				)
			})

			It("returns the error", func() {
				_, _, err := client.GetApplicationProcesses("some-app-guid")
				Expect(err).To(MatchError(ccerror.ResourceNotFoundError{Message: "App not found"}))
			})
		})
	})

	Describe("GetApplicationProcessByType", func() {
		var (
			process  Process
			warnings []string
			err      error
		)

		JustBeforeEach(func() {
			process, warnings, err = client.GetApplicationProcessByType("some-app-guid", "some-type")
		})

		Context("when the process exists", func() {
			BeforeEach(func() {
				response := `{
					"guid": "process-1-guid",
					"type": "some-type",
					"memory_in_mb": 32,
					"health_check": {
						"type": "http",
						"data": {
							"timeout": 90,
							"endpoint": "/health"
						}
					}
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes/some-type"),
						RespondWith(http.StatusOK, response, http.Header{"X-Cf-Warnings": {"this is a warning"}}),
					),
				)
			})

			It("returns the process and all warnings", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(warnings).To(ConsistOf("this is a warning"))
				Expect(process).To(Equal(Process{
					GUID:       "process-1-guid",
					Type:       "some-type",
					MemoryInMB: 32,
					HealthCheck: ProcessHealthCheck{
						Type: "http",
						Data: ProcessHealthCheckData{Endpoint: "/health"}},
				}))
			})
		})

		Context("when the application does not exist", func() {
			BeforeEach(func() {
				response := `{
					"errors": [
						{
							"detail": "Application not found",
							"title": "CF-ResourceNotFound",
							"code": 10010
						}
					]
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes/some-type"),
						RespondWith(http.StatusNotFound, response, http.Header{"X-Cf-Warnings": {"this is a warning"}}),
					),
				)
			})

			It("returns a ResourceNotFoundError", func() {
				Expect(warnings).To(ConsistOf("this is a warning"))
				Expect(err).To(MatchError(ccerror.ResourceNotFoundError{Message: "Application not found"}))
			})
		})

		Context("when the cloud controller returns errors and warnings", func() {
			BeforeEach(func() {
				response := `{
					"errors": [
						{
							"code": 10008,
							"detail": "The request is semantically invalid: command presence",
							"title": "CF-UnprocessableEntity"
						},
						{
							"code": 10009,
							"detail": "Some CC Error",
							"title": "CF-SomeNewError"
						}
					]
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v3/apps/some-app-guid/processes/some-type"),
						RespondWith(http.StatusTeapot, response, http.Header{"X-Cf-Warnings": {"this is a warning"}}),
					),
				)
			})

			It("returns the error and all warnings", func() {
				Expect(err).To(MatchError(ccerror.V3UnexpectedResponseError{
					ResponseCode: http.StatusTeapot,
					V3ErrorResponse: ccerror.V3ErrorResponse{
						Errors: []ccerror.V3Error{
							{
								Code:   10008,
								Detail: "The request is semantically invalid: command presence",
								Title:  "CF-UnprocessableEntity",
							},
							{
								Code:   10009,
								Detail: "Some CC Error",
								Title:  "CF-SomeNewError",
							},
						},
					},
				}))
				Expect(warnings).To(ConsistOf("this is a warning"))
			})
		})
	})

	Describe("PatchApplicationProcessHealthCheck", func() {
		var (
			endpoint string

			warnings []string
			err      error
		)

		JustBeforeEach(func() {
			warnings, err = client.PatchApplicationProcessHealthCheck("some-process-guid", "some-type", endpoint)
		})

		Context("when patching the process succeeds", func() {
			Context("and the endpoint is non-empty", func() {
				BeforeEach(func() {
					endpoint = "some-endpoint"
					expectedBody := `{
					"health_check": {
						"type": "some-type",
						"data": {
							"endpoint": "some-endpoint"
						}
					}
				}`
					server.AppendHandlers(
						CombineHandlers(
							VerifyRequest(http.MethodPatch, "/v3/processes/some-process-guid"),
							VerifyJSON(expectedBody),
							RespondWith(http.StatusOK, "", http.Header{"X-Cf-Warnings": {"this is a warning"}}),
						),
					)
				})

				It("patches this process's health check", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(warnings).To(ConsistOf("this is a warning"))
				})
			})

			Context("and the endpoint is empty", func() {
				BeforeEach(func() {
					endpoint = ""
					expectedBody := `{
					"health_check": {
						"type": "some-type",
						"data": {
							"endpoint": null
						}
					}
				}`
					server.AppendHandlers(
						CombineHandlers(
							VerifyRequest(http.MethodPatch, "/v3/processes/some-process-guid"),
							VerifyJSON(expectedBody),
							RespondWith(http.StatusOK, "", http.Header{"X-Cf-Warnings": {"this is a warning"}}),
						),
					)
				})

				It("patches this process's health check", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(warnings).To(ConsistOf("this is a warning"))
				})
			})
		})

		Context("when the process does not exist", func() {
			BeforeEach(func() {
				endpoint = "some-endpoint"
				response := `{
					"errors": [
						{
							"detail": "Process not found",
							"title": "CF-ResourceNotFound",
							"code": 10010
						}
					]
				}`

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodPatch, "/v3/processes/some-process-guid"),
						RespondWith(http.StatusNotFound, response, http.Header{"X-Cf-Warnings": {"this is a warning"}}),
					),
				)
			})

			It("returns an error and warnings", func() {
				Expect(err).To(MatchError(ccerror.ProcessNotFoundError{}))
				Expect(warnings).To(ConsistOf("this is a warning"))
			})
		})

		Context("when the cloud controller returns errors and warnings", func() {
			BeforeEach(func() {
				endpoint = "some-endpoint"
				response := `{
						"errors": [
							{
								"code": 10008,
								"detail": "The request is semantically invalid: command presence",
								"title": "CF-UnprocessableEntity"
							},
							{
								"code": 10009,
								"detail": "Some CC Error",
								"title": "CF-SomeNewError"
							}
						]
					}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodPatch, "/v3/processes/some-process-guid"),
						RespondWith(http.StatusTeapot, response, http.Header{"X-Cf-Warnings": {"this is a warning"}}),
					),
				)
			})

			It("returns the error and all warnings", func() {
				Expect(err).To(MatchError(ccerror.V3UnexpectedResponseError{
					ResponseCode: http.StatusTeapot,
					V3ErrorResponse: ccerror.V3ErrorResponse{
						Errors: []ccerror.V3Error{
							{
								Code:   10008,
								Detail: "The request is semantically invalid: command presence",
								Title:  "CF-UnprocessableEntity",
							},
							{
								Code:   10009,
								Detail: "Some CC Error",
								Title:  "CF-SomeNewError",
							},
						},
					},
				}))
				Expect(warnings).To(ConsistOf("this is a warning"))
			})
		})
	})
})
