package integration_test

import (
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("client", func() {
	It("repositories", func() {
		r, err := client.Repositories()
		Expect(err).ToNot(HaveOccurred())
		Expect(r).ToNot(BeEmpty())
	})

	It("can retrieve the current user", func() {
		u, err := client.CurrentUser()
		Expect(err).ToNot(HaveOccurred())
		Expect(u).ToNot(BeEmpty())
	})

	It("can retrieve a list of users", func() {
		u, err := client.Users()
		Expect(err).ToNot(HaveOccurred())
		Expect(u).ToNot(BeEmpty())

		cu, err := client.CurrentUser()
		Expect(err).ToNot(HaveOccurred())
		Expect(u).To(ContainElement(WithTransform(func(bu bitbucket.User) string { return bu.GetName() }, Equal(cu))))
	})
})
