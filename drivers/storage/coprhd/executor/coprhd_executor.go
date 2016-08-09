package executor

import (
	"github.com/akutz/gofig"

	"github.com/emccode/libstorage/api/registry"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/drivers/storage/coprhd"
	"github.com/emccode/libstorage/drivers/storage/coprhd/executor/client"
)

// driver is the storage executor for the vbox storage driver.
type driver struct {
	config gofig.Config
	client types.StorageExecutor
}

func init() {
	registry.RegisterStorageExecutor(coprhd.Name, newDriver)
}

func newDriver() types.StorageExecutor {
	return &driver{}
}

func (d *driver) Init(ctx types.Context, config gofig.Config) error {
	d.config = config
	d.client = client.New()
	return nil
}

func (d *driver) Name() string {
	return coprhd.Name
}

func (d *driver) InstanceID(
	ctx types.Context,
	opts types.Store) (*types.InstanceID, error) {
	return d.client.InstanceID(ctx, opts)
}

func (d *driver) NextDevice(
	ctx types.Context,
	opts types.Store) (string, error) {
	return "", types.ErrNotImplemented
}

func (d *driver) LocalDevices(
	ctx types.Context,
	opts *types.LocalDevicesOpts) (*types.LocalDevices, error) {
	return d.client.LocalDevices(ctx, opts)
}
