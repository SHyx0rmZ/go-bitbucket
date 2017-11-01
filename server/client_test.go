package server_test

import (
	"context"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("", func() {
	var (
		client     bitbucket.Client
		testServer *ghttp.Server
	)

	BeforeEach(func() {
		testServer = ghttp.NewServer()
		testServer.Writer = GinkgoWriter
		client, _ = server.NewClient(context.TODO(), http.DefaultClient, testServer.URL())
	})

	AfterEach(func() {
		testServer.Reset()
		testServer.Close()
	})

	It("", func() {
		Expect(client).ToNot(BeNil())
	})

	Describe("", func() {
		It("returns no projects from an empty result", func() {
			testServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/rest/api/1.0/projects"),
					ghttp.RespondWith(200, `{"isLastPage":true,"nextPageStart":null,"values":[]}`),
				),
			)

			//testServer.HandleFunc("/rest/api/1.0/projects", func(w http.ResponseWriter, r *http.Request) {
			//	w.Write([]byte(`{"isLastPage":true,"nextPageStart":null,"values":[]}`))
			//})

			p, err := client.Projects()

			Expect(err).To(BeNil())
			Expect(p).To(Equal([]bitbucket.Project{}))
		})
	})

	It("returns projects from a non-empty result", func() {
		testServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/rest/api/1.0/projects"),
				ghttp.RespondWith(200, `{"isLastPage":false,"nextPageStart":1,"values":[{}]}`),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/rest/api/1.0/projects", "start=1"),
				ghttp.RespondWith(200, `{"isLastPage":true,"nextPageStart":null,"values":[{}]}`),
			),
		)

		p, err := client.Projects()

		Expect(err).To(BeNil())
		Expect(p).ToNot(BeEmpty())
		Expect(len(p)).To(Equal(2))
	})
})
