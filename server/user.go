package server

type user struct {
	client *client

	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

func (u *user) SetClient(c *client) {
	u.client = c
}

func (u *user) GetName() string {
	return u.Name
}
