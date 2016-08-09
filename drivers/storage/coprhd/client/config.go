package client

//CoprHDClientConfig ...
type CoprHDClientConfig struct {
	block    bool
	endpoint string
	file     bool
	insecure bool
	password string
	project  string
	token    string
	username string
	varray   string
	vpool    string
}

//NewClientConfig ...
func NewClientConfig(
	block bool,
	endpoint string,
	file bool,
	insecure bool,
	password string,
	project string,
	token string,
	username string,
	varray string,
	vpool string,
) *CoprHDClientConfig {

	c = new(Config)
	c.Init()

	return c
}

// Init ...
func (c *CoprHDClientConfig) Init() *CoprHDClientConfig {
	return c
}

// Endpoint ...
func (c *CoprHDClientConfig) Endpoint() string {
	return c.endpoint
}

// Insecure ...
func (c *CoprHDClientConfig) Insecure() bool {
	return c.insecure
}

// Password ...
func (c *CoprHDClientConfig) Password() string {
	return c.password
}

// Token ...
func (c *CoprHDClientConfig) Token() string {
	return c.token
}

// Project ...
func (c *CoprHDClientConfig) Project() string {
	return c.project
}

// Block ...
func (c *CoprHDClientConfig) Block() string {
	return c.block
}

// File ...
func (c *CoprHDClientConfig) File() string {
	return c.file
}

//Username ...
func (c *CoprHDClientConfig) Username() string {
	return c.username
}

//VArray ...
func (c *CoprHDClientConfig) VArray() string {
	return c.varray
}

//VPool ...
func (c *CoprHDClientConfig) VPool() string {
	return c.varray
}
