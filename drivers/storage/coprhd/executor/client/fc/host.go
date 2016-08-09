package fc

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/akutz/goof"
)

const (
	//FCHOSTDIR ...
	FCHOSTDIR = "/sys/class/fc_host"
	//SCSIHOSTDIR ...
	SCSIHOSTDIR = "/sys/class/scsi_host"
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

//Hosts list of fc_host
func Hosts() []string {
	hosts, _ := filepath.Glob(FCHOSTDIR + "/host*")
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
	return FCHOSTDIR + "/" + c.HostID()
}

//ScsiHostPath ...
func (c *Host) ScsiHostPath() string {
	return SCSIHOSTDIR + "/" + c.HostID()
}

//NodeName ...
func (c *Host) NodeName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/node_name")
	if err != nil {
		return "", goof.Newf("->Host->NodeName(): %v", err)
	}

	return string(file), nil
}

//PortName ...
func (c *Host) PortName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_name")
	if err != nil {
		return "", goof.Newf("->Host->PortName(): %v", err)
	}

	return string(file), nil
}

//PortType ...
func (c *Host) PortType() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_type")
	if err != nil {
		return "", goof.Newf("->Host->PortType(): %v", err)
	}

	return string(file), nil
}

//FabricName ...
func (c *Host) FabricName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/fabric_name")
	if err != nil {
		return "", goof.Newf("->Host->FabricName(): %v", err)
	}

	return string(file), nil
}

//PortState ...
func (c *Host) PortState() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_state")
	if err != nil {
		return "", goof.Newf("->Host->PortState(): %v", err)
	}

	return string(file), nil
}

//PortID ...
func (c *Host) PortID() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_id")
	if err != nil {
		return "", goof.Newf("->Host->PortID(): %v", err)
	}

	return string(file), nil
}

//IssueLip ...
func (c *Host) IssueLip() error {

	file := c.BasePath() + "/issue_lip"
	if err := ioutil.WriteFile(file, []byte("1"), 0666); err != nil {
		return goof.Newf("->Host->IssueLip(): %v", err)
	}

	return nil
}

//Rescan the Specified Host
func (c *Host) Rescan() error {

	// The scan is done under <scsi_host_dir>/<hostID>/scan
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

//Sessions check for contextual sessions available and create Array of Sessions
func (c *Host) Sessions() []*Session {
	c.sessions = nil
	rports, _ := filepath.Glob(c.BasePath() + "/device/rport-[0-99]:[0-99]-[0-99]")
	for _, rport := range rports {
		c.sessions = append(c.sessions, NewSession(c, rport))
	}
	return c.sessions
}
