package storage

import (
	"fmt"
	"net"
	"strings"
	"sync"
  "crypto/tls"
  "net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"

  "github.com/go-openapi/runtime"
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/registry"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/drivers/storage/coprhd"
  coprhdcli "github.com/emccode/libstorage/drivers/storage/coprhd/client"
)

// Driver represents a vbox driver implementation of StorageDriver
type driver struct {
	sync.Mutex
	config        gofig.config
	client        *coprhdcli.Client
}

func init() {
	registry.RegisterStorageDriver(coprhd.Name, newDriver)
}

func newDriver() types.StorageDriver {
	return &driver{}
}

// Name returns the name of the driver
func (d *driver) Name() string {
	return coprhd.Name
}

// GofigEndpoint ...
func (d *driver) GofigEndpoint() string {
	return d.config.GetString("coprhd.endpoint")
}

// GofigInsecure ...
func (d *driver) GofigInsecure() bool {
	return d.config.GetBool("coprhd.insecure")
}

// GofigPassword ...
func (d *driver) GofigPassword() string {
	return d.config.GetString("coprhd.password")
}

// GofigProject ...
func (d *driver) GofigProject() string {
	return d.config.GetString("coprhd.project")
}

// GofigType ...
func (d *driver) GofigType() string {
	return d.config.GetString("coprhd.type")
}

// GofigFile ...
func (d *driver) GofigFile() string {
  if d.GofigType() == types.FILE {
    return true
  }
  return false
}

// GofigBlock ...
func (d *driver) GofigBlock() string {
  if d.GofigType() == types.BLOCK {
    return true
  }
  return false
}

// GofigUsername ...
func (d *driver) GofigUsername() string {
	return d.config.GetString("coprhd.username")
}

// GofigVArray ...
func (d *driver) GofigVArray() string {
	return d.config.GetString("coprhd.varray")
}

// GofigVPool ...
func (d *driver) GofigVPool() string {
	return d.config.GetString("coprhd.vpool")
}

// Init initializes the driver.
func (d *driver) Init(ctx types.Context, config gofig.config) (error) {
  // Set configuration GoFig
	d.config = config

  // Set log Fields
	fields := log.Fields{
		"endpoint":   d.GofigEndpoint(),
    "type":       d.GofigType(),
    "insecure":   d.GofigInsecure(),
		"password":   d.GofigPassword(),
		"project":    d.GofigProject(),
    "token":      d.GofigToken(),
		"username":   d.GofigUsername(),
		"varray":     d.GofigVArray(),
		"vpool":      d.GofigVPool(),
	}

  // Obfuscation of password in logs
	if d.GofigPassword() == "" {
		fields["password"] = ""
	} else {
		fields["password"] = "******"
	}

  // Create CoprHD clientConfig
  coprhdClientConfig := coprhdcli.NewClientConfig(
    d.GofigEndpoint(),
    d.GofigInsecure(),
    d.GofigPassword(),
    d.GofigProject(),
    d.GofigFile(),
    d.GofigBlock(),
    d.GofigUsername(),
    d.GofigVArray(),
    d.GofigVPool(),
  )

  // Create CoprHD Client with Configuration
  c.client = coprhdcli.NewClient(coprhdClientConfig)
  if err != nil {
    return goof.WithFieldsE(fields, "Error initializing driver", err)
  }

	log.WithFields(fields).Info("CoprHD Storage Driver Initialized")
	return nil
}


func (d *driver) InstanceID(ctx types.Context) (string, error) {
	return context.InstanceID(ctx)
}

// InstanceInspect returns an instance.
func (d *driver) InstanceInspect(
	ctx types.Context,
	opts types.Store) (*types.Instance, error) {

	iid := context.MustInstanceID(ctx)
	if iid.ID != "" {
		return &types.Instance{InstanceID: iid}, nil
	}

	id, err := d.InstanceID(ctx)
	if err != nil {
		return nil, err
	}

	instanceID := &types.InstanceID{ID: id, Driver: d.Name()}

	return &types.Instance{InstanceID: instanceID}, nil
}

// LocalDevices returns a map of the system's local devices.
func (d *driver) LocalDevices(
	ctx types.Context,
	opts types.Store) (*types.LocalDevices, error) {

	if ld, ok := context.LocalDevices(ctx); ok {
		return ld, nil
	}

	return nil, goof.New("missing local devices")
}

// NextDevice returns the next available device (not implemented).
func (d *driver) NextDevice(
	ctx types.Context,
	opts types.Store) (string, error) {
	return "", nil
}

// Type returns the type of storage a driver provides
// CoprHD Support Block, File, Object.
// LibStorage Driver implementation can return multiple types?
func (d *driver) Type(
  ctx types.Context) (types.StorageType, error) {
  return types.BLOCK, nil
}

// NextDeviceInfo returns the information about the driver's next available
// device workflow.
func (d *driver) NextDeviceInfo(
	ctx types.Context) (*types.NextDeviceInfo, error) {
	return nil, nil
}


// Use gocoprhd to show volume details
func (d *driver) ShowVolume(volumeID string) (*types.Volume, error) {
  showVolumeParams := gocoprhd.client.block.NewShowVolumeParams().
                                              WithID(volUrn)

  resp, err := d.client.Block.ShowVolume(showVolumeParams, d.authInfo)
  if err != nil {
    return nil, goof.WithFieldsE( fields,
                                  "error unable to show volume details",
                                  err )
  }

  // Only to be more Human readable
  volume := resp.Payload

  return &types.Volume{
    ID:               volume.ID,
    Name:             volume.Name,
    AvailabilityZone: volume.Varray.Name,
    Status:           volume.AccessState,
    Type:             types.Block,
    Size:             volume.ProvisionedCapacityGb,
    NetworkName:      volume.Wwn,
    Attachments: "",
    Fields: volume,
  }, nil

}

func (d *driver) GetVolumeExports(volumeID string) (types.VolumeAttachment, error) {
  listVolumeExportsParams := block.NewListVolumeExportsParams().
                                  WithID(volumeID)

	//use any function to do REST operations
	resp, err := d.client.Block.ListVolumeExports(listVolumeExportsParams, d.authInfo)
  if err != nil {
    return nil, goof.WithFieldsE( fields,
                                  "error unable to list volume exports",
                                  err )
  }

  // Only to be more Human readable
  exports := resp.Payload.Itl
  for _, exp := range exports {
    attachmentSD := &types.VolumeAttachment{
      // Backend or Client side ID?
      VolumeID:   volumeID,
      // Instance Reference?
      InstanceID: instanceID,
      // Backend or Client side?
      DeviceName: "",
      // Backend or Client Volume Access Status?
      Status:     volume.AccessState,
    }

  }
}

/*
Host: scsi2
Channel: 02
Id: 00
Lun: 00
Vendor: MegaRAID
Model: LD0 RAID5 34556R
Rev: 1.01
Type:   Direct-Access
ANSI SCSI revision: 02
*/

func (d *driver) getDiskInfo() {
  f, err := os.Open("/proc/scsi/scsi")
  if err != nil {
    return nil, err
  }
  defer f.Close()

  return parFileInfo(f)
}

func (d *driver)  parseFileInfo(
  r io.Reader,
  parse string  ) ( []map[string]string, error ) {
  kva := []map[string]string{}
  s := bufio.NewScanner(r)
  for i := 0; s.Scan(); i++ {
		if err := s.Err(); err != nil {
			return nil, err
		}
    _, err := fmt.Sscanf( s.Text(),
                          parse,
                          &key,
                          &value)
    if err != nil {
      return nil,  err
    }
    kva[i][key] = value
  }
  return kva
}


// List Exported Volumes
func (d *driver) Volumes(
  ctx types.Context,
	opts *types.VolumesOpts) ([]*types.Volume, error) {

	var ( volumesSD []*types.Volume
      	attachmentsSD []*types.VolumeAttachment  )


  for _, volUrn := range d.GetVolumes().ID {
    for _, volume := range d.ShowVolume(volUrn)

	  instanceID := &types.InstanceID{
	    // Global ID
			ID:     volUrn,
			Driver: d.Name(),
		}

		for _, attachment := range volumeExports.Exports {
	  attachmentSD := &types.VolumeAttachment{
	    // Backend or Client side ID?
			VolumeID:   volume.NativeID,
	    // Instance Reference
			InstanceID: instanceID,
	    // Backend or Client side
			DeviceName: volume.DeviceLabel,
	    // Backend or Client Volume Access Status
			Status:     volume.AccessState,
		}
		attachmentsSD = append(attachmentsSD, attachmentSD)
	}

	volumeSD := &types.Volume{
		Name:             volume.DeviceLabel,
		ID:               volume.NativeID,
		AvailabilityZone: d.CoprhdVArray.Name(),
		Status:           volume.AccessState,
		Type:             volume.Protocols,
		IOPS:             "",
		Size:             volume.ProvisionedCapacity,
		Attachments:      attachmentsSD,
	}
	volumesSD = append(volumesSD, volumeSD)

	return volumesSD, nil
}

// Return Volume Information if it exists
func (d *driver) VolumeInspect(
	ctx types.Context,
	volumeID string,
	opts *types.VolumeInspectOpts) (*types.Volume, error) {

	var volumeSD []*types.Volume
	var volume *gocoprhd.Volume

  volume, err := d.CoprhdVolume().Id(volumeID).Query()

/* CoprHD Volume Structure Details
	// Volume is a complete coprhd volume object
	Volume struct {
		StorageObject       `json:",inline"`
		WWN                 string      `json:"wwn"`
		Protocols           []string    `json:"protocols"`
		Protection          interface{} `json:"protection"`
		ConsistencyGroup    string      `json:"consistency_group,omitempty"`
		StorageController   string      `json:"storage_controller"`
		DeviceLabel         string      `json:"device_label"`
		NativeId            string      `json:"native_id"`
		ProvisionedCapacity string      `json:"provisioned_capacity_gb"`
		AllocatedCapacity   string      `json:"allocated_capacity_gb"`
		RequestedCapacity   string      `json:"requested_capacity_gb"`
		PreAllocationSize   string      `json:"pre_allocation_size_gb"`
		IsComposite         bool        `json:"is_composite"`
		ThinlyProvisioned   bool        `json:"thinly_provisioned"`
		HABackingVolumes    []string    `json:"high_availability_backing_volumes"`
		AccessState         string      `json:"access_state"`
		StoragePool         Resource    `json:"storage_pool"`
		Exports							[]VolumeExports `json:"exports"`
	}

// VolumeExports is the export list of the volume
	VolumeExports struct {
		Device							string	`json:"device"`
		Target							string	`json:"target"`
		Initiator						string	`json:"initiator"`
		SanZoneName					string	`json:"san_zone_name"`
		Export							VolumeExportsExport	`json:"export"`
	}

// VolumeExportsExport contains the details about the Export
	VolumeExportsExport struct {
		Id									string	`json:"id"`
		Name								string	`json:"name"`
		Link								string	`json:"link"`
	}

*/

  if err != nil {
    return nil, err
  }

	return &types.Volume{
		Name:             volume.DeviceLabel,
		ID:               volume.NativeID,
		AvailabilityZone: d.CoprhdVArray.Name(),
		Status:           volume.AccessState,
		Type:             volume.Protocols,
		IOPS:             "",
		Size:             volume.ProvisionedCapacity,
		Attachments:      attachmentsSD,
	}, nil

}

// VolumeCreate creates a new volume.
func (d *driver) VolumeCreate(
  ctx types.Context,
  volumeName string,
	opts *types.VolumeCreateOpts) (*types.Volume, error) {

  // Volume Name

	vol, err := d.Volume()
          .Name(volumeName)
          .Create(opts.Size)
	if err != nil {
		return nil, goof.WithFieldE(
      "volumeName",
      volumeName,
      "Error creating volume",
      err)
	}

	return d.VolumeInspect(
    ctx,
    volumeName,
    &types.VolumeInspectOpts{Attachments: false})
}

// VolumeRemove removes a volume.
func (d *driver) VolumeRemove(
	ctx types.Context,
	volumeID string,
	opts types.Store) error {


  err := d.Volume()
      .Name(volumeID)
      .Query()
      .Delete()
	if err != nil {
		return err
	}

	return nil
}

// VolumeAttach attaches a volume.
func (d *driver) VolumeAttach(
	ctx types.Context,
	volumeID string,
	opts *types.VolumeAttachOpts) (*types.Volume, string, error) {

	d.Lock()
	defer d.Unlock()

  d.Volume.Name(volumeID)
  d.Export.Query()
  d.Export.Initiators("string")

  export, err := d.Export.Create()
  if err != nil {
    return nil, nil, err
  }
}

// VolumeDetach detaches a volume.
func (d *driver) VolumeDetach(
	ctx types.Context,
	volumeID string,
	opts *types.VolumeDetachOpts) (*types.Volume, error) {

	d.Lock()
	defer d.Unlock()

  d.Export.Delete()
}

// VolumeCreateFromSnapshot (not implemented).
func (d *driver) VolumeCreateFromSnapshot(
	ctx types.Context,
	snapshotID, volumeName string,
	opts *types.VolumeCreateOpts) (*types.Volume, error) {
	return nil, types.ErrNotImplemented
}

// VolumeCopy copies an existing volume (not implemented)
func (d *driver) VolumeCopy(
	ctx types.Context,
	volumeID, volumeName string,
	opts types.Store) (*types.Volume, error) {
	return nil, types.ErrNotImplemented
}

// VolumeSnapshot snapshots a volume (not implemented)
func (d *driver) VolumeSnapshot(
	ctx types.Context,
	volumeID, snapshotName string,
	opts types.Store) (*types.Snapshot, error) {
	return nil, types.ErrNotImplemented
}

func (d *driver) VolumeDetachAll(
	ctx types.Context,
	volumeID string,
	opts types.Store) error {
	return nil
}

func (d *driver) Snapshots(
	ctx types.Context,
	opts types.Store) ([]*types.Snapshot, error) {
	return nil, nil
}

func (d *driver) SnapshotInspect(
	ctx types.Context,
	snapshotID string,
	opts types.Store) (*types.Snapshot, error) {
	return nil, nil
}

func (d *driver) SnapshotCopy(
	ctx types.Context,
	snapshotID, snapshotName, destinationID string,
	opts types.Store) (*types.Snapshot, error) {
	return nil, nil
}

func (d *driver) SnapshotRemove(
	ctx types.Context,
	snapshotID string,
	opts types.Store) error {

	return nil
}

func (d *driver) getVolume(ctx types.Context, volumeID, volumeName string,
	attachments bool) ([]*types.Volume, error) {
	var volumes []isi.Volume
	if volumeID != "" || volumeName != "" {
		volume, err := d.client.GetVolume(volumeID, volumeName)
		if err != nil && !strings.Contains(err.Error(), "Unable to open object") {
			return nil, err
		}
		if volume != nil {
			volumes = append(volumes, volume)
		}
	} else {
		var err error
		volumes, err = d.client.GetVolumes()
		if err != nil {
			return nil, err
		}
	}

	if len(volumes) == 0 {
		return nil, nil
	}

	var atts []*types.VolumeAttachment
	if attachments {
		var err error
		atts, err = d.getVolumeAttachments(ctx)
		if err != nil {
			return nil, err
		}
	}

	attMap := make(map[string][]*types.VolumeAttachment)
	for _, att := range atts {
		if attMap[att.VolumeID] == nil {
			attMap[att.VolumeID] = make([]*types.VolumeAttachment, 0)
		}
		attMap[att.VolumeID] = append(attMap[att.VolumeID], att)
	}

	var volumesSD []*types.Volume
	for _, volume := range volumes {
		volSize, err := d.getSize(volume.Name, volume.Name)
		if err != nil {
			return nil, err
		}

		vatts, _ := attMap[volume.Name]
		volumeSD := &types.Volume{
			Name:        volume.Name,
			ID:          volume.Name,
			Size:        volSize,
			Attachments: vatts,
		}
		volumesSD = append(volumesSD, volumeSD)
	}

	return volumesSD, nil
}
