package integration_test

import (
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("client", func() {
	It("can retrieve a list of projects", func() {
		projects, err := client.Projects()
		Expect(err).ToNot(HaveOccurred())
		Expect(projects).ToNot(BeEmpty())
	})

	It("can create a project", func() {
		project, err := client.CreateProject()
		Expect(err).To(SatisfyAny(
			Not(HaveOccurred()),
			WithTransform(func(e *bitbucket.Error) string { return e.ExceptionName() }, Equal("com.atlassian.bitbucket.AuthorisationException")),
		))

		if err == nil {
			client.DeleteProject(project.GetKey())
		} else {
			Skip("User is missing PROJECT_CREATE permission")
		}
	})

	It("can delete a project", func() {
		project, err := client.CreateProject()
		Expect(err).To(SatisfyAny(
			Not(HaveOccurred()),
			WithTransform(func(e *bitbucket.Error) string { return e.ExceptionName() }, Equal("com.atlassian.bitbucket.AuthorisationException")),
		))

		if err == nil {
			err = client.DeleteProject(project.GetKey())
			Expect(err).ToNot(HaveOccurred())
		} else {
			Skip("User is missing PROJECT_CREATE permission")
		}
	})
})
