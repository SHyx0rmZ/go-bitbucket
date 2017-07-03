package server

import (
	"fmt"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
)

type project struct {
	client *client

	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
}

func (p *project) SetClient(c *client) {
	p.client = c
}

func (p project) GetKey() string {
	return p.Key
}

func (p *project) Repositories() ([]bitbucket.Repository, error) {
	repositories := make([]repository, 0, 0)

	fmt.Printf("%#v\n", p)
	err := p.client.pagedRequest("/rest/api/1.0/projects/"+p.Key+"/repos", &repositories)
	if err != nil {
		return nil, err
	}

	bitbucketRepositories := make([]bitbucket.Repository, 0, len(repositories))

	for _, repository := range repositories {
		//if ca, ok := repository.Project.(clientAware); ok {
		//	ca.SetClient(p.client)
		//}

		bitbucketRepositories = append(bitbucketRepositories, &repository)
	}

	return bitbucketRepositories, nil
}
