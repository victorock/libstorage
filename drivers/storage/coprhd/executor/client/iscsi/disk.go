package iscsi

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/akutz/goof"
)

const (
	//DISKDEVDIR ...
	DISKDEVDIR = "/dev"
	//DISKBYIDDIR ...
	DISKBYIDDIR = "/dev/disk/by-id"
	//DISKBYPATHDIR ...
	DISKBYPATHDIR = "/dev/disk/by-path"
	//DISKBYUUIDDIR ...
	DISKBYUUIDDIR = "/dev/disk/by-uuid"
)

// Disk ...
type Disk struct {
	// Reference to the upstream Target Object
	lun  *Lun
	name string
}

// NewDisk is Generic Constructor
func NewDisk(lun *Lun, name string) *Disk {
	c := new(Disk)
	c.Init(lun, name)
	return c
}

//Disks list disks from all hosts, sessions, targets and luns
// ls -1d /sys/devices/platform/host*/session*/target*\:*\:*/*\:*\:*\:*/block/sd*
func Disks() []string {
	disks, _ := filepath.Glob(ISCSIHOSTDIR + "/host*/session*/target*:*:*/*:*:*:*/block/sd*")
	return disks
}

// Init ...
func (c *Disk) Init(lun *Lun, name string) *Disk {
	c.SetDiskName(name).
		SetLun(lun)
	return c
}

//SetLun ...
func (c *Disk) SetLun(lun *Lun) *Disk {
	c.lun = lun
	return c
}

//Lun ...
func (c *Disk) Lun() *Lun {
	return c.lun
}

//SetDiskName ...
func (c *Disk) SetDiskName(name string) *Disk {
	c.name = filepath.Dir(name)
	return c
}

//DiskName ...
func (c *Disk) DiskName() string {
	return c.name
}

//DevPath ...
func (c *Disk) DevPath() string {
	return DISKDEVDIR + "/" + c.DiskName()
}

//BasePath ...
func (c *Disk) BasePath() string {
	return c.Lun().BasePath() + "/" + c.DiskName()
}

//DiskByID ...
func (c *Disk) DiskByID() string {
	disksByID, _ := filepath.Glob(DISKBYIDDIR + "/wwn*")
	for _, diskID := range disksByID {
		devPath, _ := filepath.EvalSymlinks(fmt.Sprintf("%s/%s", DISKBYIDDIR, diskID))
		if devPath == c.DevPath() {
			return diskID
		}
	}
	return ""
}

//DiskByUUID ...
func (c *Disk) DiskByUUID() string {
	disksByUUID, _ := filepath.Glob(DISKBYUUIDDIR + "/*")
	for _, diskUUID := range disksByUUID {
		devPath, _ := filepath.EvalSymlinks(fmt.Sprintf("%s/%s", DISKBYUUIDDIR, diskUUID))
		if devPath == c.DevPath() {
			return diskUUID
		}
	}
	return ""
}

//DiskByPath ...
func (c *Disk) DiskByPath() string {
	disksByPath, _ := filepath.Glob(DISKBYPATHDIR + "/*")
	for _, diskPath := range disksByPath {
		devPath, _ := filepath.EvalSymlinks(fmt.Sprintf("%s/%s", DISKBYPATHDIR, diskPath))
		if devPath == c.DevPath() {
			return diskPath
		}
	}
	return ""
}

//Size ...
func (c *Disk) Size() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/size")
	if err != nil {
		return "", goof.Newf("->Disk->Size(): %v", err)
	}

	return string(file), nil
}

//Removable ...
func (c *Disk) Removable() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/removable")
	if err != nil {
		return "", goof.Newf("->Disk->Removable(): %v", err)
	}

	return string(file), nil
}

//ReadOnly ...
func (c *Disk) ReadOnly() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/ro")
	if err != nil {
		return "", goof.Newf("->Disk->ReadOnly(): %v", err)
	}

	return string(file), nil
}

//Capability ...
func (c *Disk) Capability() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/capability")
	if err != nil {
		return "", goof.Newf("->Disk->Capability(): %v", err)
	}

	return string(file), nil
}

//Dev ...
func (c *Disk) Dev() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/dev")
	if err != nil {
		return "", goof.Newf("->Disk->Dev(): %v", err)
	}

	return string(file), nil
}
