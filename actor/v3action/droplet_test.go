package v3action_test

import (
	"errors"
	"net/url"

	. "code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/actor/v3action/v3actionfakes"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Droplet Actions", func() {
	var (
		actor                     *Actor
		fakeCloudControllerClient *v3actionfakes.FakeCloudControllerClient
	)

	BeforeEach(func() {
		fakeCloudControllerClient = new(v3actionfakes.FakeCloudControllerClient)
		actor = NewActor(fakeCloudControllerClient, nil)
	})

	Describe("SetApplicationDroplet", func() {
		Context("when there are no client errors", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetApplicationsReturns(
					[]ccv3.Application{
						{GUID: "some-app-guid"},
					},
					[]string{"get-applications-warning"},
					nil,
				)

				fakeCloudControllerClient.SetApplicationDropletReturns(
					ccv3.Relationship{GUID: "some-droplet-guid"},
					[]string{"set-application-droplet-warning"},
					nil,
				)
			})

			It("sets the app's droplet", func() {
				warnings, err := actor.SetApplicationDroplet("some-app-name", "some-space-guid", "some-droplet-guid")

				Expect(err).ToNot(HaveOccurred())
				Expect(warnings).To(ConsistOf("get-applications-warning", "set-application-droplet-warning"))

				Expect(fakeCloudControllerClient.GetApplicationsCallCount()).To(Equal(1))
				queryURL := fakeCloudControllerClient.GetApplicationsArgsForCall(0)
				query := url.Values{"names": []string{"some-app-name"}, "space_guids": []string{"some-space-guid"}}
				Expect(queryURL).To(Equal(query))

				Expect(fakeCloudControllerClient.SetApplicationDropletCallCount()).To(Equal(1))
				appGUID, dropletGUID := fakeCloudControllerClient.SetApplicationDropletArgsForCall(0)
				Expect(appGUID).To(Equal("some-app-guid"))
				Expect(dropletGUID).To(Equal("some-droplet-guid"))
			})
		})

		Context("when getting the application fails", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("some get application error")

				fakeCloudControllerClient.GetApplicationsReturns(
					[]ccv3.Application{},
					[]string{"get-applications-warning"},
					expectedErr,
				)
			})

			It("returns the error", func() {
				warnings, err := actor.SetApplicationDroplet("some-app-name", "some-space-guid", "some-droplet-guid")

				Expect(err).To(Equal(expectedErr))
				Expect(warnings).To(ConsistOf("get-applications-warning"))
			})
		})

		Context("when setting the droplet fails", func() {
			var expectedErr error
			BeforeEach(func() {
				expectedErr = errors.New("some set application-droplet error")
				fakeCloudControllerClient.GetApplicationsReturns(
					[]ccv3.Application{
						{GUID: "some-app-guid"},
					},
					[]string{"get-applications-warning"},
					nil,
				)

				fakeCloudControllerClient.SetApplicationDropletReturns(
					ccv3.Relationship{},
					[]string{"set-application-droplet-warning"},
					expectedErr,
				)
			})

			It("returns the error", func() {
				warnings, err := actor.SetApplicationDroplet("some-app-name", "some-space-guid", "some-droplet-guid")

				Expect(err).To(Equal(expectedErr))
				Expect(warnings).To(ConsistOf("get-applications-warning", "set-application-droplet-warning"))
			})
		})
	})
})
