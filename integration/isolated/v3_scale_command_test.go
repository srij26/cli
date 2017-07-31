package isolated

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("v3-scale command", func() {
	var (
		orgName   string
		spaceName string
		appName   string
		userName  string
	)

	BeforeEach(func() {
		orgName = helpers.NewOrgName()
		spaceName = helpers.NewSpaceName()
		appName = helpers.PrefixedRandomName("app")
		userName, _ = helpers.GetCredentials()
	})

	Describe("help", func() {
		Context("when --help flag is set", func() {
			It("displays command usage to output", func() {
				session := helpers.CF("v3-scale", "--help")

				Eventually(session.Out).Should(Say("NAME:"))
				Eventually(session.Out).Should(Say("v3-scale - Change or view the instance count, disk space limit, and memory limit for an app"))

				Eventually(session.Out).Should(Say("USAGE:"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]"))

				Eventually(session.Out).Should(Say("OPTIONS:"))
				Eventually(session.Out).Should(Say("-i\\s+Number of instances"))
				Eventually(session.Out).Should(Say("-k\\s+Disk limit \\(e\\.g\\. 256M, 1024M, 1G\\)"))
				Eventually(session.Out).Should(Say("-m\\s+Memory limit \\(e\\.g\\. 256M, 1024M, 1G\\)"))

				Eventually(session).Should(Exit(0))
			})
		})
	})

	Context("when the environment is not setup correctly", func() {
		Context("when no API endpoint is set", func() {
			BeforeEach(func() {
				helpers.UnsetAPI()
			})

			It("fails with no API endpoint set message", func() {
				session := helpers.CF("v3-scale", appName)
				Eventually(session).Should(Say("FAILED"))
				Eventually(session.Err).Should(Say("No API endpoint set\\. Use 'cf login' or 'cf api' to target an endpoint\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when not logged in", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
			})

			It("fails with not logged in message", func() {
				session := helpers.CF("v3-scale", appName)
				Eventually(session).Should(Say("FAILED"))
				Eventually(session.Err).Should(Say("Not logged in\\. Use 'cf login' to log in\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when there is no org set", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
				helpers.LoginCF()
			})

			It("fails with no org targeted error message", func() {
				session := helpers.CF("v3-scale", appName)
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Err).Should(Say("No org targeted, use 'cf target -o ORG' to target an org\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when there is no space set", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
				helpers.LoginCF()
				helpers.TargetOrg(ReadOnlyOrg)
			})

			It("fails with no space targeted error message", func() {
				session := helpers.CF("v3-scale", appName)
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Err).Should(Say("No space targeted, use 'cf target -s SPACE' to target a space\\."))
				Eventually(session).Should(Exit(1))
			})
		})
	})

	Context("when the environment is set up correctly", func() {
		BeforeEach(func() {
			setupCF(orgName, spaceName)
		})

		Context("when the app name is not provided", func() {
			It("tells the user that the app name is required, prints help text, and exits 1", func() {
				session := helpers.CF("v3-scale")

				Eventually(session.Err).Should(Say("Incorrect Usage: the required argument `APP_NAME` was not provided"))
				Eventually(session.Out).Should(Say("NAME:"))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when the app does not exist", func() {
			It("displays app not found and exits 1", func() {
				invalidAppName := "invalid-app-name"
				session := helpers.CF("v3-scale", invalidAppName)
				Eventually(session.Out).Should(Say("Showing health and status for app %s in org %s / space %s as %s\\.\\.\\.", invalidAppName, orgName, spaceName, userName))
				Eventually(session.Err).Should(Say("App %s not found", invalidAppName))
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when the app exists", func() {
			BeforeEach(func() {
				helpers.WithHelloWorldApp(func(appDir string) {
					Eventually(helpers.CustomCF(helpers.CFEnv{WorkingDirectory: appDir}, "v3-push", appName)).Should(Exit(0))
				})
			})

			Context("when flag options are not provided", func() {
			  It("displays the current scale properties", func() {
					session := helpers.CF("v3-scale", appName)
					Eventually(session.Out).Should(Say("Showing current scale of app %s in org %s / space %s as %s\\.\\.\\.", appName, orgName, spaceName, userName))
					Eventually(session.Out).Should(Say("memory: 32M"))
					Eventually(session.Out).Should(Say("disk: 1G"))
					Eventually(session.Out).Should(Say("instances: 1"))
					Eventually(session).Should(Exit(0))
			  })
			})

			Context("when flag options are provided", func() {
				It("scales the app accordingly", func() {
					session := helpers.CF("v3-scale", appName)
					Eventually(session.Out).Should(Say("memory: 32M"))
					Eventually(session.Out).Should(Say("disk: 1G"))
					Eventually(session.Out).Should(Say("instances: 1"))
					Eventually(session).Should(Exit(0))

					session = helpers.CF("v3-scale", appName, "-i", "3")
					Eventually(session.Out).Should(Say("Scaling app %s in org %s / space %s as %s\\.\\.\\.", appName, orgName, spaceName, userName))
					Eventually(session).Should(Exit(0))

					session = helpers.CF("v3-scale", appName)
					Eventually(session.Out).Should(Say("memory: 32M"))
					Eventually(session.Out).Should(Say("disk: 1G"))
					Eventually(session.Out).Should(Say("instances: 3"))
					Eventually(session).Should(Exit(0))

					session = helpers.CF("v3-scale", appName, "-k", "92M", "-m", "64M")
					Eventually(session.Out).Should(Say("Scaling app %s in org %s / space %s as %s\\.\\.\\.", appName, orgName, spaceName, userName))
					Eventually(session).Should(Exit(0))

					session = helpers.CF("v3-scale", appName)
					Eventually(session.Out).Should(Say("memory: 64M"))
					Eventually(session.Out).Should(Say("disk: 92M"))
					Eventually(session.Out).Should(Say("instances: 3"))
					Eventually(session).Should(Exit(0))
				})
			})
		})
	})

	Context("when invalid scale option values are provided", func() {
		Context("when a negative value is passed to a flag argument", func() {
			It("outputs an error message to the user, provides help text, and exits 1", func() {
				session := helpers.CF("v3-scale", "some-app", "-i=-5")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-i' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))

				session = helpers.CF("v3-scale", "some-app", "-k=-5")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-k' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))

				session = helpers.CF("v3-scale", "some-app", "-m=-5")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-m' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when a non-integer value is passed to a flag argument", func() {
			It("outputs an error message to the user, provides help text, and exits 1", func() {
				session := helpers.CF("v3-scale", "some-app", "-i", "not-an-integer")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-i' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))

				session = helpers.CF("v3-scale", "some-app", "-k", "not-an-integer")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-k' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))

				session = helpers.CF("v3-scale", "some-app", "-m", "not-an-integer")
				Eventually(session.Err).Should(Say("Incorrect Usage: invalid argument for flag `-m' \\(expected int\\)"))
				Eventually(session.Out).Should(Say("cf v3-scale APP_NAME \\[-i INSTANCES\\] \\[-k DISK\\] \\[-m MEMORY\\]")) // help
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when the unit of measurement is not provided", func() {
			It("outputs an error message to the user, provides help text, and exits 1", func() {
				session := helpers.CF("v3-scale", "some-app", "-k", "9")
				Eventually(session.Err).Should(Say("Invalid memory quota: 9"))
				Eventually(session.Err).Should(Say("Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB"))
				Eventually(session).Should(Exit(1))

				session = helpers.CF("v3-scale", "some-app", "-m", "7")
				Eventually(session.Err).Should(Say("Invalid memory quota: 7"))
				Eventually(session.Err).Should(Say("Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB"))
				Eventually(session).Should(Exit(1))
			})
		})
	})
})
