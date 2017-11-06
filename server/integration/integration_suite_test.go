package integration_test

import (
	"context"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
	"testing"
)

func TestServerIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
}

var endpoint = os.Getenv("BITBUCKET_SERVER_TESTING_ENDPOINT")
var username = os.Getenv("BITBUCKET_SERVER_TESTING_USERNAME")
var password = os.Getenv("BITBUCKET_SERVER_TESTING_PASSWORD")

var client bitbucket.Client

var _ = BeforeSuite(func() {
	if endpoint != "" {
		Ω(endpoint).ShouldNot(BeEmpty(), "must specify $BITBUCKET_SERVER_TESTING_ENDPOINT")
		Ω(username).ShouldNot(BeEmpty(), "must specify $BITBUCKET_SERVER_TESTING_USERNAME")
		Ω(password).ShouldNot(BeEmpty(), "must specify $BITBUCKET_SERVER_TESTING_PASSWORD")

		var err error
		ctx := context.WithValue(context.Background(), bitbucket.BitbucketAuth, &bitbucket.BasicAuth{
			Username: username,
			Password: password,
		})
		client, err = server.NewClient(ctx, http.DefaultClient, endpoint)
		Expect(err).ToNot(HaveOccurred())
	}
})

var _ = BeforeEach(func() {
	if client == nil {
		Skip("Environment variables need to be set for Bitbucket Server integration")
	}
})

func TestServerIntegrationForReal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitbucket Server Integration Suite")
}
