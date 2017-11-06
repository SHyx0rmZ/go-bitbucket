package cloud_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//bitbucket "github.com/SHyx0rmZ/go-bitbucket/bitbucket"
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

var _ = Describe("client", func() {
	var (
		client     *cloud.Client
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

	It("can retrieve the current user", func() {
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

	It("can not retrieve a list of users", func() {
		resp, err := ioutil.ReadFile("testdata/2.0/users/get_200.json")
		Expect(err).ToNot(HaveOccurred())

		testServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/2.0/users"),
				ghttp.RespondWith(200, resp),
			),
		)

		users, err := client.Users()
		Expect(err).ToNot(HaveOccurred())
		Expect(users).To(BeEmpty())
	})

	It("can retrieve a list of teams", func() {
		resp, err := ioutil.ReadFile("testdata/2.0/teams/get_200.json")
		Expect(err).ToNot(HaveOccurred())

		testServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/2.0/teams"),
				ghttp.RespondWith(200, resp),
			),
		)

		teams, err := client.Teams()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(teams)).To(Equal(1))
		Expect(teams[0].GetName()).To(Equal("some-team"))
	})

	It("can retrieve a list of team members", func() {
		respTeams, err := ioutil.ReadFile("testdata/2.0/teams/get_200.json")
		Expect(err).ToNot(HaveOccurred())

		respMembers, err := ioutil.ReadFile("testdata/2.0/teams/some-team/members/get_200.json")
		Expect(err).ToNot(HaveOccurred())

		testServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/2.0/teams"),
				ghttp.RespondWith(200, respTeams),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/2.0/teams/some-team/members"),
				ghttp.RespondWith(200, respMembers),
			),
		)

		teams, err := client.Teams()
		Expect(err).ToNot(HaveOccurred())

		members, err := teams[0].Members()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(members)).To(Equal(1))
		Expect(members[0].GetName()).To(Equal("some-user"))
	})
})
