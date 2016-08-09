package client

import (
	"crypto/tls"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/goof"
	runtime "github.com/go-openapi/runtime"
	httpclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/victorock/gocoprhd/client"
)

//CoprHDClient ...
type CoprHDClient struct {
	config    *CoprHDClientConfig
	transport *httpclient.Runtime
	authInfo  runtime.ClientAuthInfoWriter
	client    *apiclient.CoprHD
}

//NewClient ...
func NewClient(config *CoprHDClientConfig) (*CoprHDClient, error) {
	// Create object
	c := new(CoprHDClient)

	//Create Config Object
	c.config = config

	// Initialize the client
	c.Init()

	log.Info("CoprHD Client: Created")
	return c, nil
}

// Init ...
func (c *CoprHDClient) Init() (*CoprHDClient, error) {

	// Init Endpoint()
	if _, err := c.Endpoint(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Insecure()
	if _, err := c.Insecure(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Authentication()
	if _, err := c.Authentication(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Token()
	if _, err := c.Token(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Login()
	if _, err := c.Login(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// create the API client, with the transport
	c.client = apiclient.New(c.transport, strfmt.Default)

	log.Info("CoprHD Client: Initialized")
	return c, nil
}

// Endpoint Create transport Object as per user Endpoint configuration
func (c *CoprHDClient) Endpoint() (*CoprHDClient, error) {
	if c.config.Endpoint() == "" {
		return nil, goof.New("->Endpoint(): Endpoint if not configured.")
	}
	// Set the driver Endpoint
	c.transport = httpclient.New(c.config.Endpoint(), "/", []string{"https"})
	return c, nil
}

// Insecure Set transport objects insecure as per user Insecure configuration
func (c *CoprHDClient) Insecure() (*CoprHDClient, error) {

	// Set insecure transport
	if c.config.Insecure() {
		c.transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return c, nil
}

// Authentication with username/password as per user Username/Password configuration
func (c *CoprHDClient) Authentication() (*CoprHDClient, error) {
	if c.config.Username() == "" || c.config.Password() == "" {
		return nil, goof.New("->Authentication(): Username/Password is Empty")
	}

	// Set Authentication Info
	authInfo := httpclient.BasicAuth(c.config.Username(), c.config.Password())

	// Set Pointer
	c.authInfo = authInfo

	return c, nil
}

//Token Authentication with Token as per user Token configuration
func (c *CoprHDClient) Token() (*CoprHDClient, error) {

	// Initialize the Driver Token Header
	if c.config.Token() != "" {
		authInfo := httpclient.APIKeyAuth("X-SDS-AUTH-TOKEN", "header", c.config.Token())
		// Populate the Header with the token from now on
		c.authInfo = authInfo
	}

	return c, nil
}

// Login to CoprHD
func (c *CoprHDClient) Login() (*CoprHDClient, error) {

	// Initialize the Driver Login Method
	login, err := c.client.Authentication.Login(nil, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->Login(), %v", err)
	}

	// Populate the Header with the token from now on
	c.authInfo = httpclient.APIKeyAuth("X-SDS-AUTH-TOKEN", "header", login.XSDSAUTHTOKEN)
	log.Info("CoprHD Client: Login()")
	return c, nil
}

// TaskCheck to CoprHD
func (c *CoprHDClient) TaskCheck() (*CoprHDClient, error) {
	return c, nil
}

// Volumes Use gocoprhd to get the list of volumes
func (c *CoprHDClient) Volumes() ([]string, error) {
	resp, err := c.client.Block.ListVolumes(nil, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->ListVolumes(), %v", err)
	}
	return resp.Payload.ID, nil
}
