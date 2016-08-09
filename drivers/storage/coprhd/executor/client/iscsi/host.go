package iscsi

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/akutz/goof"
)

const (
	//ISCSIHOSTDIR ...
	ISCSIHOSTDIR = "/sys/devices/platform"
)

// Host ...
type Host struct {
	Session
	Disk
	Lun
	Target
	ID       string
	sessions []*Session
}

// New is Generic Constructor
//func NewHost(machine Machine, ID string) *Host {
func New(ID string) *Host {
	c := new(Host)
	c.Init(ID)
	return c
}

//Hosts ... list all iscsi_host
func Hosts() []string {
	hosts, _ := filepath.Glob(ISCSIHOSTDIR + "/host*")
	return hosts
}

// Init ...
//func (c *Host) Init(ID string) *Host {
func (c *Host) Init(ID string) *Host {
	c.SetHostID(ID)
	return c
}

//HostID ...
func (c *Host) HostID() string {
	return c.ID
}

//SetHostID ...
func (c *Host) SetHostID(ID string) *Host {
	c.ID = filepath.Dir(ID)
	return c
}

//BasePath is the base Path...
func (c *Host) BasePath() string {
	return ISCSIHOSTDIR + "/" + c.HostID()
}

//ScsiHostPath ...
func (c *Host) ScsiHostPath() string {
	return c.BasePath() + "/device/scsi_host/" + c.HostID()
}

//InitiatorName ...
func (c *Host) InitiatorName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/initiatorname")
	if err != nil {
		return "", goof.Newf("->Host->InitiatorName(): %v", err)
	}

	return string(file), nil
}

//IPAddress ...
func (c *Host) IPAddress() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/ipaddress")
	if err != nil {
		return "", goof.Newf("->Host->IPAddress(): %v", err)
	}

	return string(file), nil
}

//HWAddress ...
func (c *Host) HWAddress() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/hwaddress")
	if err != nil {
		return "", goof.Newf("->Host->HWAddress(): %v", err)
	}

	return string(file), nil
}

//NetDev ...
func (c *Host) NetDev() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/netdev")
	if err != nil {
		return "", goof.Newf("->Host->NetDev(): %v", err)
	}

	return string(file), nil

}

//Rescan the Specified Host
func (c *Host) Rescan() error {

	// The scan is done under <Path>/scsi_host/<hostID>/scan
	file, _ := filepath.Glob(c.ScsiHostPath() + "/scan")
	if file == nil {
		return goof.Newf("->Host->Rescan(): File Access Error, %v/scan", c.ScsiHostPath())
	}

	scanFile := file[0]
	if err := ioutil.WriteFile(scanFile, []byte("- - -"), 0666); err != nil {
		return goof.Newf("->Host->Rescan(): File Write Error, %v/scan", c.ScsiHostPath())
	}

	time.Sleep(2 * time.Second)
	return nil
}

//Sessions check for sessions available and create Array of Sessions
func (c *Host) Sessions() []*Session {
	c.sessions = nil
	IDs, _ := filepath.Glob(c.BasePath() + "/session[1-99]")
	for _, ID := range IDs {
		c.sessions = append(c.sessions, NewSession(c, ID))
	}
	return c.sessions
}
