package cloud

import (
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"net/http"
	"errors"
	"io"
	"strings"
	"encoding/json"
	"bytes"
)


type client struct {
	httpClient *http.Client
	endpoint string
	auth bitbucket.Auth
}

func NewClient() (bitbucket.Client, error) {
	return &client{
		httpClient: http.DefaultClient,
		endpoint:   "https://api.bitbucket.org/2.0/",
	}, nil
}

func (c *client) do(method string, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, strings.TrimRight(c.endpoint, "/")+url, body)
	if err != nil {
		return nil, err
	}

	if basicAuth, ok := c.auth.(*bitbucket.BasicAuth); ok {
		request.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}

	if method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *client) request(url string, v interface{}) error {
	response, err := c.do("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) requestPost(url string, v interface{}, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := c.do("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		return errors.New(response.Status)
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Users() ([]bitbucket.User, error) {
	var user user

	err := c.request("user", &user)
	if err != nil {
		return nil, err
	}

	return []*user {&user}, nil
}

func (client) Projects() ([]bitbucket.Project, error) {
	return nil, errors.New("Not Implemented")
}

func (client) Repository(path string) (bitbucket.Repository, error) {
	return nil, errors.New("Not Implemented")
}
