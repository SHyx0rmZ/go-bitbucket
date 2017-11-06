package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type client struct {
	endpoint string
	client   *http.Client
	auth     bitbucket.Auth
}

type clientAware interface {
	SetClient(c *client)
}

func NewClient(ctx context.Context, hc *http.Client, endpoint string) (bitbucket.Client, error) {
	if hc == nil {
		hc = contextHTTPClient(ctx)
	}
	auth := contextBitbucketAuth(ctx)

	return &client{
		endpoint: endpoint,
		client:   hc,
		auth:     auth,
	}, nil
}

func contextBitbucketAuth(ctx context.Context) bitbucket.Auth {
	if ctx != nil {
		if a, ok := ctx.Value(bitbucket.BitbucketAuth).(bitbucket.Auth); ok {
			return a
		}
	}
	return nil
}

func contextHTTPClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(bitbucket.HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	return http.DefaultClient
}

func (c *client) Projects() ([]bitbucket.Project, error) {
	projects := make([]project, 0, 0)

	err := c.pagedRequest("/rest/api/1.0/projects", &projects)
	if err != nil {
		return nil, err
	}

	bitbucketProjects := make([]bitbucket.Project, len(projects))

	for index := range projects {
		bitbucketProjects[index] = &projects[index]
	}

	return bitbucketProjects, nil
}

func (c *client) Repository(name string) (bitbucket.Repository, error) {
	var url string

	components := strings.Split(name, "/")

	if len(components) != 2 {
		return nil, errors.New("Invalid repository: " + name)
	}

	ownerName := components[0]
	repositoryName := components[1]

	if strings.Index(ownerName, "~") == 0 {
		url = "/rest/api/1.0/users/" + strings.TrimLeft(ownerName, "~") + "/repos/" + repositoryName
	} else {
		url = "/rest/api/1.0/projects/" + ownerName + "/repos/" + repositoryName
	}

	var result repository

	err := c.request(url, &result)
	if err != nil {
		return nil, err
	}

	result.SetClient(c)

	return &result, nil
}

func (c *client) Repositories() ([]bitbucket.Repository, error) {
	repositories := make([]repository, 0, 0)

	err := c.pagedRequest("/rest/api/1.0/repos?permission=REPO_WRITE", &repositories)
	if err != nil {
		return nil, err
	}

	bitbucketRepositories := make([]bitbucket.Repository, len(repositories))
	for index := range repositories {
		bitbucketRepositories[index] = &repositories[index]
	}

	return bitbucketRepositories, nil
}

func (c *client) CreateRepository(path string) (bitbucket.Repository, error) {
	return nil, errors.New("not yet implemented")
}

func (c *client) CurrentUser() (string, error) {
	response, err := c.do("GET", "/rest/api/1.0/users?limit=0", strings.NewReader(""))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
	}

	return response.Header.Get("X-Ausername"), nil
}

func (c *client) Users() ([]bitbucket.User, error) {
	users := make([]user, 0, 0)

	err := c.pagedRequest("/rest/api/1.0/users", &users)
	if err != nil {
		return nil, err
	}

	bitbucketUsers := make([]bitbucket.User, len(users))
	for index := range users {
		bitbucketUsers[index] = &users[index]
	}

	return bitbucketUsers, nil
}

func (c *client) SetHTTPClient(hc *http.Client) {
	c.client = hc
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

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	if ca, ok := v.(*clientAware); ok {
		(*ca).SetClient(c)
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

	defer response.Body.Close()

	if response.StatusCode != 201 {
		return errors.New(response.Status)
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	if ca, ok := v.(*clientAware); ok {
		(*ca).SetClient(c)
	}

	return nil
}

type PagedResult struct {
	IsLastPage    bool              `json:"isLastPage"`
	Values        []json.RawMessage `json:"values,omitempty"`
	NextPageStart *int              `json:"nextPageStart"`
}

func (c *client) pagedRequest(url string, v interface{}) error {
	resultValue := reflect.ValueOf(v)

	if resultValue.Kind() != reflect.Ptr || resultValue.IsNil() {
		return errors.New("Invalid return type")
	}

	resultList := reflect.ValueOf(v).Elem()
	resultElemType := resultList.Type().Elem()

	var pageStart *int = nil

	for {
		fullUrl := url

		if pageStart != nil {
			if strings.Contains(fullUrl, "?") {
				fullUrl += "&start=" + strconv.Itoa(*pageStart)
			} else {
				fullUrl += "?start=" + strconv.Itoa(*pageStart)
			}
		}

		var results PagedResult

		err := c.request(fullUrl, &results)
		if err != nil {
			return err
		}

		for _, jsonBytes := range results.Values {
			newResult := reflect.New(resultElemType).Elem()

			err = json.Unmarshal(jsonBytes, newResult.Addr().Interface())
			if err != nil {
				return err
			}

			if ca, ok := newResult.Addr().Interface().(clientAware); ok {
				ca.SetClient(c)
			}

			resultList.Set(reflect.Append(resultList, newResult))
		}

		if results.IsLastPage == true || results.NextPageStart == nil {
			break
		}

		pageStart = results.NextPageStart
	}

	return nil
}
