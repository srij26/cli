package isolated

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/integration/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("login command", func() {
	var buffer *Buffer

	FContext("when the API endpoint is not set", func() {
		BeforeEach(func() {
			username, password := helpers.GetCredentials()
			helpers.UnsetAPI()
			buffer = NewBuffer()
			buffer.Write([]byte("api.bosh-lite.com\n"))
			buffer.Write([]byte(fmt.Sprintf("%s\n", username)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
			buffer.Write([]byte("\n"))
			fmt.Fprintf(os.Stderr, "Username = %s\n", username)
			fmt.Fprintf(os.Stderr, "Password = %s\n", password)
		})

		It("prompts the user for an endpoint", func() {
			session := helpers.CFWithStdin(buffer, "login")
			Eventually(session.Out).Should(Say("API endpoint>"))
			Eventually(session.Out).Should(Say("Email>"))
			Eventually(session.Out).Should(Say("Password>"))

			Eventually(session).Should(Exit(0))
		})
	})

	Context("when no arguments are passed in", func() {
		BeforeEach(func() {
			username, password := helpers.GetCredentials()
			buffer = NewBuffer()
			buffer.Write([]byte(fmt.Sprintf("%s\n", username)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
		})

		It("prompts the user for email, password, org, and space", func() {
			session := helpers.CFWithStdin(buffer, "login")
			Eventually(session.Out).Should(Say("Email>"))
			Eventually(session.Out).Should(Say("Password>"))

			Eventually(session).Should(Exit(0))
		})
	})

	// need to create a dedicated user and org for:
	// when only one org exists
	// - with no spaces
	// - with one space
	// - with two spaces
	// when there are multiple orgs

	Describe("tests which need to create a user", func() {
		var (
			username string
			password string
		)

		BeforeEach(func() {
			helpers.LoginCF()
			username = helpers.RandomUsername()
			password = helpers.RandomPassword()
			Expect(helpers.CF("create-user", username, password)).NotTo(HaveOccurred())
			buffer = NewBuffer()
			buffer.Write([]byte(fmt.Sprintf("%s\n", username)))
			buffer.Write([]byte(fmt.Sprintf("%s\n", password)))
		})

		Context("when the user sees no orgs", func() {
			It("prompts for username and password but does not target any org", func() {
			})
		})

		Context("when the user sees one org", func() {
			var orgName string

			BeforeEach(func() {
				orgName = helpers.NewOrgName()

				Expect(helpers.CF("create-org", orgName))
				Expect(helpers.CF("set-org-role", username, orgName, "OrgManager"))
			})

			Context("when the user can see no spaces", func() {
				It("prompts for username and password then automatically targets the existing org but no space", func() {
				})
			})

			Context("when the user can see one space", func() {
				BeforeEach(func() {
					spaceName := helpers.NewOrgName()

					Expect(helpers.CF("create-space", spaceName, "-o", orgName))
					Expect(helpers.CF("set-space-role", username, orgName, spaceName, "SpaceManager"))
				})

				It("prompts for username and password then automatically targets the existing org and space", func() {
				})
			})

			Context("when the user can see multiple spaces", func() {
				var spaceName string

				BeforeEach(func() {
					spaceName = helpers.NewOrgName()

					Expect(helpers.CF("create-space", spaceName, "-o", orgName))
					Expect(helpers.CF("set-space-role", username, orgName, spaceName, "SpaceManager"))
				})

				It("prompts for username and password then automatically targets the existing org and prompts for space", func() {
				})
			})
		})

		Context("when the user sees multiple orgs", func() {
			var orgName string

			BeforeEach(func() {
				orgName = helpers.NewOrgName()

				Expect(helpers.CF("create-org", orgName))
				Expect(helpers.CF("set-org-role", username, orgName, "OrgManager"))
				otherOrgName := helpers.NewOrgName()

				Expect(helpers.CF("create-org", otherOrgName))
				Expect(helpers.CF("set-org-role", username, otherOrgName, "OrgManager"))
			})

			Context("when the user can see no spaces", func() {
				It("prompts for username and password then prompts the user for an existing org but no space", func() {
				})
			})

			Context("when the user can see one space", func() {
				BeforeEach(func() {
					spaceName := helpers.NewOrgName()

					Expect(helpers.CF("create-space", spaceName, "-o", orgName))
					Expect(helpers.CF("set-space-role", username, orgName, spaceName, "SpaceManager"))
				})

				It("prompts for username and password then prompts the user for an existing org and targets its sole space", func() {
				})
			})

			Context("when the user can see multiple spaces", func() {
				var spaceName string

				BeforeEach(func() {
					spaceName = helpers.NewOrgName()

					Expect(helpers.CF("create-space", spaceName, "-o", orgName))
					Expect(helpers.CF("set-space-role", username, orgName, spaceName, "SpaceManager"))

					otherSpaceName := helpers.NewOrgName()

					Expect(helpers.CF("create-space", otherSpaceName, "-o", orgName))
					Expect(helpers.CF("set-space-role", username, orgName, otherSpaceName, "SpaceManager"))
				})

				It("prompts for username and password then automatically prompts the user for an existing org and space", func() {
				})
			})
		})
	})

	// Context("when the API endpoint is set", func() {
	// 	Context("when there are no arguments", func() {
	// 		BeforeEach(func() {
	// 			helpers.UnsetAPI()
	// 			buffer = NewBuffer()
	// 			buffer.Write([]byte("\n"))
	// 		})

	// 		It("prompts the user for an email address and password", func() {
	// 		})
	// 	})
	// })

	Context("when --sso-passcode flag is given", func() {
		Context("when a passcode isn't provided", func() {
			It("prompts the user to try again", func() {
				session := helpers.CFWithStdin(buffer, "login", "--sso-passcode")
				Eventually(session.Err).Should(Say("Incorrect Usage: expected argument for flag `--sso-passcode'"))
			})
		})

		Context("when the provided passcode is invalid", func() {
			It("prompts the user to try again", func() {
				session := helpers.CFWithStdin(buffer, "login", "--sso-passcode", "bad-passcode")
				Eventually(session.Out).Should(Say("Authenticating..."))
				Eventually(session.Out).Should(Say("Credentials were rejected, please try again."))
			})
		})
	})

	Context("when both --sso and --sso-passcode flags are provided", func() {
		It("errors with invalid use", func() {
			session := helpers.CFWithStdin(buffer, "login", "--sso", "--sso-passcode", "some-passcode")
			Eventually(session.Out).Should(Say("Incorrect usage: --sso-passcode flag cannot be used with --sso"))
			Eventually(session).Should(Exit(1))
		})
	})
})
