package fc

import (
	"io/ioutil"
	"path/filepath"

	"github.com/akutz/goof"
)

const (
	//FCSESSIONDIR ...
	FCSESSIONDIR = "/sys/class/fc_remote_ports"
)

// Session ...
type Session struct {
	// Reference for Upstream host Object
	host    *Host
	ID      string
	targets []*Target
}

// NewSession is Generic Constructor
func NewSession(host *Host, ID string) *Session {
	c := new(Session)
	c.Init(host, ID)
	return c
}

//Sessions list luns from all hosts, sessions and targets
// ls -1d /sys/class/fc_remote_ports/rport-*:*-*
func Sessions() []string {
	sessions, _ := filepath.Glob(FCSESSIONDIR + "/rport-*:*-*")
	return sessions
}

// Init ...
func (c *Session) Init(host *Host, ID string) *Session {
	c.SetHost(host).
		SetSessionID(ID)
	return c
}

//SessionID ...
func (c *Session) SessionID() string {
	return c.ID
}

//SetSessionID ...
func (c *Session) SetSessionID(ID string) *Session {
	c.ID = filepath.Dir(ID)
	return c
}

//SetHost reference for upstream object calling session
func (c *Session) SetHost(host *Host) *Session {
	c.host = host
	return c
}

//Host is reference to the Host object Calling Session
func (c *Session) Host() *Host {
	return c.host
}

//BasePath ...
func (c *Session) BasePath() string {
	return FCSESSIONDIR + "/" + c.SessionID()
}

//PortName ...
func (c *Session) PortName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_name")
	if err != nil {
		return "", goof.Newf("->Session->PortName()%v", err)
	}

	return string(file), nil
}

//PortID ...
func (c *Session) PortID() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_id")
	if err != nil {
		return "", goof.Newf("->Session->PortID()%v", err)
	}

	return string(file), nil
}

//PortState ...
func (c *Session) PortState() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/port_state")
	if err != nil {
		return "", goof.Newf("->Session->PortState()%v", err)
	}

	return string(file), nil
}

//NodeName ...
func (c *Session) NodeName() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/node_name")
	if err != nil {
		return "", goof.Newf("->Session->NodeName()%v", err)
	}

	return string(file), nil
}

//DevLossTmo ...
func (c *Session) DevLossTmo() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/dev_loss_tmo")
	if err != nil {
		return "", goof.Newf("->Session->DevLossTmo()%v", err)
	}

	return string(file), nil
}

//SupportedClasses ...
func (c *Session) SupportedClasses() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/supported_classes")
	if err != nil {
		return "", goof.Newf("->Session->SupportedClasses()%v", err)
	}

	return string(file), nil
}

//MaxFrameSize ...
func (c *Session) MaxFrameSize() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/maxframe_size")
	if err != nil {
		return "", goof.Newf("->Session->MaxFrameSize()%v", err)
	}

	return string(file), nil
}

//Targets contextual this session item.
func (c *Session) Targets() []*Target {
	c.targets = nil
	//target[Host]:[Channel]:[Id]
	IDs, _ := filepath.Glob(c.BasePath() + "/device/target[0-9]:[0-9]:[0-9]")
	for _, ID := range IDs {
		c.targets = append(c.targets, NewTarget(c, ID))
	}
	return c.targets
}
