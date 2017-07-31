package v3_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/translatableerror"
	"code.cloudfoundry.org/cli/command/v3"
	"code.cloudfoundry.org/cli/command/v3/v3fakes"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = FDescribe("Scale Command", func() {
	var (
		cmd             v3.V3ScaleCommand
		testUI          *ui.UI
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		fakeActor       *v3fakes.FakeV3ScaleActor
		binaryName      string
		executeErr      error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)
		fakeActor = new(v3fakes.FakeV3ScaleActor)

		cmd = v3.V3ScaleCommand{
			UI:          testUI,
			Config:      fakeConfig,
			SharedActor: fakeSharedActor,
			Actor:       fakeActor,
		}

		cmd.RequiredArgs.AppName = "some-app"

		binaryName = "faceman"
		fakeConfig.BinaryNameReturns(binaryName)
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("when checking target fails", func() {
		BeforeEach(func() {
			fakeSharedActor.CheckTargetReturns(sharedaction.NotLoggedInError{BinaryName: binaryName})
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError(translatableerror.NotLoggedInError{BinaryName: binaryName}))

			Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
			_, checkTargetedOrg, checkTargetedSpace := fakeSharedActor.CheckTargetArgsForCall(0)
			Expect(checkTargetedOrg).To(BeTrue())
			Expect(checkTargetedSpace).To(BeTrue())
		})
	})

	Context("when the user is logged in, and org and space are targeted", func() {
		BeforeEach(func() {
			fakeConfig.HasTargetedOrganizationReturns(true)
			fakeConfig.TargetedOrganizationReturns(configv3.Organization{Name: "some-org"})
			fakeConfig.HasTargetedSpaceReturns(true)
			fakeConfig.TargetedSpaceReturns(configv3.Space{
				GUID: "some-space-guid",
				Name: "some-space"})
			fakeConfig.CurrentUserReturns(
				configv3.User{Name: "some-user"},
				nil)
		})

		Context("when getting the current user returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("getting current user error")
				fakeConfig.CurrentUserReturns(
					configv3.User{},
					expectedErr)
			})

			It("returns the error", func() {
				Expect(executeErr).To(MatchError(expectedErr))
			})
		})

		Context("when the application does not exist", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = v3action.ApplicationNotFoundError{Name: "some-app"}
				fakeActor.GetProcessByApplicationNameAndSpaceReturns(v3action.Process{}, v3action.Warnings{"warning-1", "warning-2"}, expectedErr)
			})

			It("returns an ApplicationNotFoundError and all warnings", func() {
				Expect(executeErr).To(Equal(translatableerror.ApplicationNotFoundError{Name: "some-app"}))
				Expect(testUI.Out).To(Say("Showing current scale of app some-app in org some-org / space some-space as some-user\\.\\.\\."))

				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))
			})
		})

		FContext("when the application exists", func() {
			Context("when no flag options are provided", func() {
				BeforeEach(func() {
					instance1 := ccv3.Instance{
						Index: 1,
					}
					instance2 := ccv3.Instance{
						Index: 2,
					}
					process := v3action.Process{
						Type: "web",
						Instances: []v3action.Instance{
							v3action.Instance(instance1),
							v3action.Instance(instance2),
						},
						MemoryInMB: 128,
						DiskInMB:   2000,
					}

					fakeActor.GetProcessByApplicationNameAndSpaceReturns(process, v3action.Warnings{"warning-1", "warning-2"}, nil)
				})

				It("displays current scale properties and all warnings", func() {
					Expect(executeErr).ToNot(HaveOccurred())
					Expect(testUI.Out).ToNot(Say("Scaling"))
					Expect(testUI.Out).To(Say("Showing current scale of app some-app in org some-org / space some-space as some-user\\.\\.\\."))
					Expect(testUI.Out).To(Say("memory:\\s+128M"))
					Expect(testUI.Out).To(Say("disk:\\s+2G"))
					Expect(testUI.Out).To(Say("instances:\\s+2"))
					Expect(testUI.Out).To(Say("OK"))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))

					Expect(fakeActor.ScaleProcessByApplicationNameAndSpaceCallCount()).To(Equal(0))
					Expect(fakeActor.GetProcessByApplicationNameAndSpaceCallCount()).To(Equal(1))
					appName, spaceGUID := fakeActor.GetProcessByApplicationNameAndSpaceArgsForCall(0)
					Expect(appName).To(Equal("some-app"))
					Expect(spaceGUID).To(Equal("some-space-guid"))
				})

				Context("when an error is encountered getting process information", func() {
					var expectedErr error

					BeforeEach(func() {
						expectedErr = errors.New("get process error")
						fakeActor.GetProcessByApplicationNameAndSpaceReturns(
							v3action.Process{},
							v3action.Warnings{"get-process-warning"},
							expectedErr,
						)
					})

					It("returns the error and displays all warnings", func() {
						Expect(executeErr).To(Equal(expectedErr))
						Expect(testUI.Err).To(Say("get-process-warning"))
					})
				})
			})

			Context("when all flag options are provided", func() {
				BeforeEach(func() {
					cmd.DiskLimit = flag.Megabytes{Size: 64}
					cmd.MemoryLimit = flag.Megabytes{Size: 256}
					cmd.Instances = 3

					fakeActor.ScaleProcessByApplicationNameAndSpaceReturns(
						v3action.Warnings{"scale-warning-1", "scale-warning-2"},
						nil,
					)
				})

				It("scales the application", func() {
					Expect(executeErr).ToNot(HaveOccurred())
					Expect(testUI.Out).ToNot(Say("Showing current"))
					Expect(testUI.Out).To(Say("Scaling app some-app in org some-org / space some-space as some-user\\.\\.\\."))
					Expect(testUI.Out).To(Say("OK"))

					Expect(testUI.Err).To(Say("scale-warning-1"))
					Expect(testUI.Err).To(Say("scale-warning-2"))

					Expect(fakeActor.GetProcessByApplicationNameAndSpaceCallCount()).To(Equal(0))
					Expect(fakeActor.ScaleProcessByApplicationNameAndSpaceCallCount()).To(Equal(1))
					appName, spaceGUID, process := fakeActor.ScaleProcessByApplicationNameAndSpaceArgsForCall(0)
					Expect(appName).To(Equal("some-app"))
					Expect(spaceGUID).To(Equal("some-space-guid"))
					Expect(process).To(Equal(ccv3.Process{
						Type:       "web",
						Instances:  3,
						DiskInMB:   64,
						MemoryInMB: 256,
					}))
				})

				Context("when an error is encountered scaling the application", func() {
					var expectedErr error

					BeforeEach(func() {
						expectedErr = errors.New("scale process error")
						fakeActor.ScaleProcessByApplicationNameAndSpaceReturns(
							v3action.Warnings{"scale-process-warning"},
							expectedErr,
						)
					})

					It("returns the error and displays all warnings", func() {
						Expect(executeErr).To(Equal(expectedErr))
						Expect(testUI.Err).To(Say("scale-process-warning"))
					})
				})
			})

			Context("when only the instances flag option is provided", func() {
				Context("when the instances flag is set to 0", func() {
					It("scales the application", func() {

					})
				})

				Context("when the instances flag option is set to > 0", func() {
					It("scales the application", func() {

					})
				})
			})

			Context("when only the instances flag option is not provided", func() {
				It("scales the disk and memory of the application and restarts the application", func() {

				})
			})
		})
	})
})
