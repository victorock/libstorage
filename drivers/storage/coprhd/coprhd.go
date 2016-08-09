package coprhd

import (
	"github.com/akutz/gofig"
)

const (
	// Name is the name of the storage driver
	Name = "coprhd"
)

func init() {
	registerConfig()
}

func registerConfig() {
	r := gofig.NewRegistration("CoprHD")
	r.Key(gofig.String, "", "localhost:4443", "", "coprhd.endpoint")
	r.Key(gofig.Bool, "", true, "", "coprhd.insecure")
	r.Key(gofig.String, "", "root", "", "coprhd.username")
	r.Key(gofig.String, "", "ChangeMe1!", "", "coprhd.password")
	r.Key(gofig.String, "", "", "", "coprhd.token")
	r.Key(gofig.String, "", "", "urn:storageos:Project:7d46540b-140c-4f39-91b8-52d276356cf0:global", "coprhd.project")
	r.Key(gofig.String, "", "", "urn:storageos:VirtualArray:ad18dd81-99c6-415d-9081-6091db3df599:vdc1", "coprhd.varray")
	r.Key(gofig.String, "", "", "urn:storageos:VirtualPool:7e036b4a-9cba-4357-9afc-3aa7539f10c0:vdc1", "coprhd.vpool")
	// Block || File
	r.Key(gofig.String, "", "", "block", "coprhd.type")
	r.Key(gofig.String, "", "", "3.0", "coprhd.version")
	gofig.Register(r)
}
