package cloud_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bitbucket "github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/cloud"
	"github.com/onsi/gomega/ghttp"
	"io/ioutil"
	"net/http"
)

const (
	roleAdmin       = "admin"
	roleContributor = "contributor"
	roleMember      = "member"
)

var _ = Describe("", func() {
	var (
		client     bitbucket.Client
		testServer *ghttp.Server
	)

	BeforeEach(func() {
		testServer = ghttp.NewServer()
		testServer.Writer = GinkgoWriter
		client, _ = cloud.NewClient(http.DefaultClient, testServer.URL())
	})

	AfterEach(func() {
		testServer.Reset()
		testServer.Close()
	})

	It("", func() {
		resp, err := ioutil.ReadFile("testdata/2.0/user/get_200.json")
		Expect(err).ToNot(HaveOccurred())

		testServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/2.0/user"),
				ghttp.RespondWith(200, resp),
			),
		)

		user, err := client.CurrentUser()
		Expect(err).ToNot(HaveOccurred())
		Expect(user).To(Equal("some-user"))
	})
})
