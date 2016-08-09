package iscsi

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/akutz/goof"
)

// Lun ...
type Lun struct {
	// Reference to the upstream Target Object
	target *Target
	ID     string
	disks  []*Disk
}

// NewLun is Generic Constructor
func NewLun(target *Target, ID string) *Lun {
	c := new(Lun)
	c.Init(target, ID)
	return c
}

//Luns list luns from all hosts, sessions and targets
// ls -1d /sys/devices/platform/host*/session*/target*\:*\:*/*\:*\:*\:*/
func Luns() []string {
	luns, _ := filepath.Glob(ISCSIHOSTDIR + "/host*/session*/target*:*:*/*:*:*:*")
	return luns
}

// Init ...
func (c *Lun) Init(target *Target, ID string) *Lun {
	c.SetLunID(ID).
		SetTarget(target)
	return c
}

//SetTarget ...
func (c *Lun) SetTarget(target *Target) *Lun {
	c.target = target
	return c
}

//Target ...
func (c *Lun) Target() *Target {
	return c.target
}

//SetLunID ...
func (c *Lun) SetLunID(ID string) *Lun {
	c.ID = filepath.Dir(ID)
	return c
}

//LunID ...
func (c *Lun) LunID() string {
	return c.ID
}

//BasePath get the reference of the SessionDir in the upstream Session Object
func (c *Lun) BasePath() string {
	return c.Target().BasePath() + "/" + c.LunID()
}

//Model ...
func (c *Lun) Model() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/model")
	if err != nil {
		return "", goof.Newf("->Model: %v", err)
	}

	return string(file), nil
}

//Vendor ...
func (c *Lun) Vendor() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/vendor")
	if err != nil {
		return "", goof.Newf("->Vendor: %v", err)
	}

	return string(file), nil
}

//Rev ...
func (c *Lun) Rev() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/rev")
	if err != nil {
		return "", goof.Newf("->Rev: %v", err)
	}

	return string(file), nil
}

//State ...
func (c *Lun) State() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/state")
	if err != nil {
		return "", goof.Newf("->State: %v", err)
	}

	return string(file), nil
}

//QueueDepth ...
func (c *Lun) QueueDepth() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/queue_depth")
	if err != nil {
		return "", goof.Newf("->QueueDepth: %v", err)
	}

	return string(file), nil
}

//IoerrCnt ...
func (c *Lun) IoerrCnt() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/ioerr_cnt")
	if err != nil {
		return "", goof.Newf("->IoerrCnt: %v", err)
	}

	return string(file), nil
}

//IodoneCnt ...
func (c *Lun) IodoneCnt() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/iodone_cnt")
	if err != nil {
		return "", goof.Newf("->IodoneCnt: %v", err)
	}

	return string(file), nil
}

//Rescan the Lun, ex: resize
func (c *Lun) Rescan() error {

	file, _ := filepath.Glob(c.BasePath() + "/rescan")
	if file == nil {
		return goof.Newf("->Rescan: File Access Error, %v/rescan", c.BasePath())
	}

	rescanFile := file[0]
	if err := ioutil.WriteFile(rescanFile, []byte("- - -"), 0666); err != nil {
		return goof.Newf("->Rescan: File Write Error, , %v/rescan", c.BasePath())
	}

	time.Sleep(1 * time.Second)
	return nil
}

//Disks check for sessions available and create Array of Sessions
func (c *Lun) Disks() []*Disk {
	c.disks = nil
	names, _ := filepath.Glob(c.BasePath() + "/block/sd*")
	for _, name := range names {
		c.disks = append(c.disks, NewDisk(c, name))
	}
	return c.disks
}
