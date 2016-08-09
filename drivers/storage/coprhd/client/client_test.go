package client

import (
	"os"
	"testing"
)

func TestInit() *CoprHDClient {

	endpoint := "localhost:4443"
	insecure := true
	password := "password"
	project := "urn:storageos:Project:7d46540b-140c-4f39-91b8-52d276356cf0:global"
	file := false
	block := true
	username := "root"
	varray := "urn:storageos:VirtualArray:ad18dd81-99c6-415d-9081-6091db3df599:vdc1"
	vpool := "urn:storageos:VirtualPool:7e036b4a-9cba-4357-9afc-3aa7539f10c0:vdc1"

	if os.Getenv("GOCOPRHD_ENDPOINT") != "" {
		endpoint = os.Getenv("GOCOPRHD_ENDPOINT")
	}

	if os.Getenv("GOCOPRHD_INSECURE") != "" {
		insecure = os.Getenv("GOCOPRHD_INSECURE")
	}

	if os.Getenv("GOCOPRHD_PASSWORD") != "" {
		password = os.Getenv("GOCOPRHD_PASSWORD")
	}

	if os.Getenv("GOCOPRHD_PROJECT") != "" {
		project = os.Getenv("GOCOPRHD_PROJECT")
	}

	if os.Getenv("GOCOPRHD_FILE") != "" {
		file = true
	}

	if os.Getenv("GOCOPRHD_BLOCK") != "" {
		block = true
	}

	if os.Getenv("GOCOPRHD_USERNAME") != "" {
		username = os.Getenv("GOCOPRHD_USERNAME")
	}

	if os.Getenv("GOCOPRHD_VARRAY") != "" {
		varray = os.Getenv("GOCOPRHD_VARRAY")
	}

	if os.Getenv("GOCOPRHD_VPOOL") != "" {
		vpool = os.Getenv("GOCOPRHD_VPOOL")
	}

	coprhdClientConfig := NewClientConfig(
		endpoint,
		insecure,
		password,
		project,
		file,
		block,
		username,
		varray,
		vpool,
	)

	coprhdClient := NewClient(coprhdClientConfig)

	return coprhdClient
}

func TestClient(t *testing.T) {
	cli := TestInit()

}
