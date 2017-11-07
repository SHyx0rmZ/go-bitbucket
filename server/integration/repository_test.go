package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("client", func() {
	It("can retrieve a list of repositories", func() {
		r, err := client.Repositories()
		Expect(err).ToNot(HaveOccurred())
		Expect(r).ToNot(BeEmpty())
	})
})
