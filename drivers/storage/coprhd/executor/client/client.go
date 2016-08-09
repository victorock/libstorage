package client

import (
	"strings"

	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/drivers/storage/coprhd"
	"github.com/emccode/libstorage/drivers/storage/coprhd/executor/client/fc"
	"github.com/emccode/libstorage/drivers/storage/coprhd/executor/client/iscsi"
	"github.com/emccode/libstorage/drivers/storage/coprhd/executor/client/utils"
)

//Executor ...
type Executor struct {
	fiberChannel []*fc.Host
	internetScsi []*iscsi.Host
	//scaleIO []*scaleio
	//networkFs []*nas
}

//New ...
func New() *Executor {
	c := new(Executor)
	c.Init()
	return c
}

//Init ...
func (c *Executor) Init() *Executor {
	return c
}

// InstanceID implementation of client
func (c *Executor) InstanceID(
	ctx types.Context,
	opts types.Store) (*types.InstanceID, error) {

	id, err := c.ID()
	if err != nil {
		return nil, err
	}

	iid := &types.InstanceID{Driver: coprhd.Name}
	if err := iid.MarshalMetadata(id); err != nil {
		return nil, err
	}

	return iid, nil
}

//LocalDevices ...
func (c *Executor) LocalDevices(
	ctx types.Context,
	opts *types.LocalDevicesOpts) (*types.LocalDevices, error) {

	mapDiskByID := make(map[string]string)

	if opts.ScanType == types.DeviceScanDeep {
		c.Rescan()
	}

	// FC Cards
	for _, fcHost := range c.Fc() {
		for _, disk := range fcHost.Disks() {
			if name := disk.DiskName(); name != "" {
				mapDiskByID[name] = c.ShortDeviceID(disk.DiskByID())
			}
		}
	}

	// Iscsi Sessions
	for _, fcHost := range c.Fc() {
		for _, disk := range fcHost.Disks() {
			if name := disk.DiskName(); name != "" {
				mapDiskByID[name] = c.ShortDeviceID(disk.DiskByID())
			}
		}
	}

	return &types.LocalDevices{
		Driver:    coprhd.Name,
		DeviceMap: mapDiskByID,
	}, nil
}

//ID ...
func (c *Executor) ID() (string, error) {
	// MachineID() || HostID()
	id, err := utils.MachineID()
	if err != nil {
		return "", err
	}

	return id, nil
}

//Iscsi objects are created in-time...
// We are not in control of actions taken directly at O.S
func (c *Executor) Iscsi() []*iscsi.Host {
	var internetScsi []*iscsi.Host
	for _, iscsiHost := range iscsi.Hosts() {
		internetScsi = append(internetScsi, iscsi.New(iscsiHost))
	}
	c.internetScsi = internetScsi
	return c.internetScsi
}

//Fc objects are created in-time...
// We are not in control of actions taken directly at O.S
func (c *Executor) Fc() []*fc.Host {
	var fiberChannel []*fc.Host
	for _, fcHost := range fc.Hosts() {
		fiberChannel = append(fiberChannel, fc.New(fcHost))
	}
	c.fiberChannel = fiberChannel
	return c.fiberChannel
}

// Rescan everywhere...
func (c *Executor) Rescan() error {

	//Trigger Rescan at each FC Card.
	for _, f := range c.Fc() {
		if err := f.Rescan(); err != nil {
			return err
		}
	}

	// Trigger Rescan at each Iscsi Session
	for _, i := range c.Iscsi() {
		if err := i.Rescan(); err != nil {
			return err
		}
	}

	return nil
}

//ShortDeviceID return short version of device iiscsi..
func (c *Executor) ShortDeviceID(dev string) string {
	sid := strings.Split(dev, "wwn")
	if len(sid) < 1 {
		return ""
	}
	aid := strings.Split(sid[1], "-")
	if len(aid) < 1 {
		return ""
	}
	return aid[0]
}

// TODO
// ScaleIO ...
//func ScaleIO() []scaleio {}
// NFS ...
//func NFS() []nfs {}
