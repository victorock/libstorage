package iscsi

import "path/filepath"

//Target ...
type Target struct {
	// Reference for Upstream Session Object
	session *Session
	//target[0-9]:[0-9]:[0-9]
	ID string
	//[0-9]:[0-9]:[0-9]:[0-255]
	luns []*Lun
}

// NewTarget is Generic Constructor
func NewTarget(session *Session, ID string) *Target {
	c := new(Target)
	c.Init(session, ID)
	return c
}

//Targets list targets from all hosts and sessions
// ls -1d /sys/devices/platform/host*/session*/target*\:*\:*/*\:*\:*\:*/
func Targets() []string {
	targets, _ := filepath.Glob(ISCSIHOSTDIR + "/host*/session*/target*:*:*")
	return targets
}

// Init ...
func (c *Target) Init(session *Session, ID string) *Target {
	c.SetTargetID(ID).
		SetSession(session)
	return c
}

//SetSession ...
func (c *Target) SetSession(session *Session) *Target {
	c.session = session
	return c
}

//Session ...
func (c *Target) Session() *Session {
	return c.session
}

//SetTargetID ...
func (c *Target) SetTargetID(ID string) *Target {
	c.ID = filepath.Dir(ID)
	return c
}

//TargetID ...
func (c *Target) TargetID() string {
	return c.ID
}

//BasePath get the reference of the Path in the upstream Session Object
func (c *Target) BasePath() string {
	return c.Session().BasePath() + "/device/" + c.TargetID()
}

//Luns ...
func (c *Target) Luns() []*Lun {
	c.luns = nil
	IDs, _ := filepath.Glob(c.BasePath() + "/[0-99]:[0-99]:[0-99]:[0-255]")
	for _, ID := range IDs {
		c.luns = append(c.luns, NewLun(c, ID))
	}
	return c.luns
}
