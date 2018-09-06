package api

// UserInfo basic user info class
type UserInfo struct {
	Username   string
	Password   string
	FirstName  string
	LastName   string
	MiddleName string
}

// Set sets user creds
func (c *UserInfo) Set(u, p string) {
	c.Username, c.Password = u, p
}

// Get gets user creds
func (c *UserInfo) Get() (string, string) {
	return c.Username, c.Password
}
