package client

import (
	"crypto/tls"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/goof"
	"github.com/emccode/libstorage/api/types"
	runtime "github.com/go-openapi/runtime"
	runtransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/victorock/gocoprhd/client"
	apiauth "github.com/victorock/gocoprhd/client/authentication"
	apiblock "github.com/victorock/gocoprhd/client/block"
	apicompute "github.com/victorock/gocoprhd/client/compute"
	apivdc "github.com/victorock/gocoprhd/client/vdc"
	apimodels "github.com/victorock/gocoprhd/models"
)

//CoprHDClient ...
type CoprHDClient struct {
	config    *CoprHDClientConfig
	transport *runtransport.Runtime
	authInfo  runtime.ClientAuthInfoWriter
	client    *apiclient.CoprHD
	block     *apiblock.Client
	vdc       *apivdc.Client
	compute   *apicompute.Client
	auth      *apiauth.Client
	//file    *apifile.Client
}

//NewClient ...
func NewClient(config *CoprHDClientConfig) (*CoprHDClient, error) {
	// Create object
	c := new(CoprHDClient)

	//Create Config Object
	c.config = config

	// Initialize the client
	if _, err := c.Init(); err != nil {
		return nil, goof.Newf("Unable Create CoprHD Client, %v", err)
	}

	log.Info("CoprHD Client: Created")
	return c, nil
}

// Init ...
func (c *CoprHDClient) Init() (*CoprHDClient, error) {

	// Init Endpoint()
	if _, err := c.Endpoint(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Insecure()
	if _, err := c.Insecure(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Authentication()
	if _, err := c.Authentication(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Token()
	if _, err := c.Token(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// Init Login()
	if _, err := c.Login(); err != nil {
		return nil, goof.Newf("->CoprHD Client %v", err)
	}

	// create the API client, with the transport
	c.client = apiclient.New(c.transport, strfmt.Default)
	c.auth = c.client.Authentication
	c.block = c.client.Block
	c.vdc = c.client.Vdc
	c.compute = c.client.Compute

	log.Info("CoprHD Client: Initialized")
	return c, nil
}

// Endpoint Create transport Object as per user Endpoint configuration
func (c *CoprHDClient) Endpoint() (*CoprHDClient, error) {
	if c.config.Endpoint() == "" {
		return nil, goof.New("->Endpoint(): Endpoint if not configured.")
	}
	// Set the driver Endpoint
	c.transport = runtransport.New(c.config.Endpoint(), "/", []string{"https"})
	return c, nil
}

// Insecure Set transport objects insecure as per user Insecure configuration
func (c *CoprHDClient) Insecure() (*CoprHDClient, error) {

	// Set insecure transport
	if c.config.Insecure() {
		c.transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return c, nil
}

// Authentication with username/password as per user Username/Password configuration
func (c *CoprHDClient) Authentication() (*CoprHDClient, error) {
	if c.config.Username() == "" || c.config.Password() == "" {
		return nil, goof.New("->Authentication(): Username/Password is Empty")
	}

	// Set Authentication Info
	c.authInfo = runtransport.BasicAuth(c.config.Username(), c.config.Password())

	return c, nil
}

//Token Authentication with Token as per user Token configuration
func (c *CoprHDClient) Token() (*CoprHDClient, error) {

	// Initialize the Driver Token Header
	if c.config.Token() != "" {
		authInfo := runtransport.APIKeyAuth("X-SDS-AUTH-TOKEN", "header", c.config.Token())
		// Populate the Header with the token from now on
		c.authInfo = authInfo
	}

	return c, nil
}

// Login to CoprHD
func (c *CoprHDClient) Login() (*CoprHDClient, error) {

	// Initialize the Driver Login Method
	login, err := c.auth.Login(nil, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->Login(), %v", err)
	}

	// Populate the Header with the token from now on
	c.authInfo = runtransport.APIKeyAuth("X-SDS-AUTH-TOKEN", "header", login.XSDSAUTHTOKEN)
	log.Info("CoprHD Client: Login()")
	return c, nil
}

// Task return task information
func (c *CoprHDClient) Task(taskID string) (*apimodels.Task, error) {

	// Create Object to Request
	showTaskParams := apivdc.NewShowTaskParams().WithID(taskID)

	//use any function to do REST operations
	resp, err := c.vdc.ShowTask(showTaskParams, authInfo)
	if err != nil {
		return nil, goof.Newf("->Task(), %#v", err)
	}

	log.Infof("->Task(): %#v", resp.Payload)
	return resp.Payload, nil
}

// AsyncTask to CoprHD
func (c *CoprHDClient) AsyncTask(task *apimodels.Task) (*apimodels.Task, error) {

	// This process will run and clockwise, giving tasks 10 minutes to finish
	// tic, tac...
	// TODO: Make the timeout value configurable?
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(10 * time.Minute)
		timeout <- true
	}()

	// We create two channels to monitor error and status of the task...
	successCh := make(chan *apimodels.Task, 1)
	errorCh := make(chan error, 1)
	go func(task *apimodels.Task) {
		log.Info("->AsyncTask(): Waiting for Task to Complete...")

		for {

			taskInfo, err := c.Task(task.ID)
			if err != nil {
				errorCh <- goof.Newf("->AsyncTask(): %#v", err)
				return
			}

			if taskInfo.Progress == 100 {
				successCh <- taskInfo
				return
			}

			time.Sleep(10 * time.Second)
		}
	}(task)

	// Monitor for messages in the channel announcing the end of task or error
	select {
	case taskDone := <-successCh:
		log.Infof("->AsyncTask(): Done, %#v", taskDone)
		return taskDone, nil
	case err := <-errorCh:
		return nil, err
	case <-timeout:
		return nil, goof.New("->AsyncTask(): Task Timeout")
	}

}

// ListVolumes Use gocoprhd to get the list of volumes
// <- Volumes
func (c *CoprHDClient) ListVolumes() (*apimodels.Volumes, error) {
	// This method doesn't requires any parameter.
	resp, err := c.block.ListVolumes(nil, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->ListVolumes(), %v", err)
	}
	return resp.Payload, nil
}

// ShowVolume Show details about the volume
// <- VolumeInspect
func (c *CoprHDClient) ShowVolume(volid string) (*apimodels.Volume, error) {

	// Create the Request Parameters
	showVolumeParams := apiblock.NewShowVolumeParams().WithID(volid)

	// Send the request
	resp, err := c.block.ShowVolume(showVolumeParams, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->ShowVolume(): %#v", err)
	}
	return resp.Payload, nil
}

// CreateVolume Use gocoprhd to create volume
// <- VolumeCreate
func (c *CoprHDClient) CreateVolume(name string, sizeGB string) (*apimodels.Volume, error) {

	// Create the request Parameters
	body := &apimodels.CreateVolume{
		ConsistencyGroup: "",
		Count:            1,
		Name:             name,
		Size:             sizeGB + "GB",
		Project:          c.config.Project(),
		Varray:           c.config.VArray(),
		Vpool:            c.config.VPool(),
	}

	// Create the New Params and Populate the Body
	CreateVolumeParams := apiblock.NewCreateVolumeParams().WithBody(body)

	//use any function to do REST operations
	resp, err := c.block.CreateVolume(CreateVolumeParams, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->CreateVolume(): %v", err)
	}

	// AsyncTask check status progress bar...
	// We always create one volume per time, but the driver supports bulk creation..
	// Therefore each volume part of the bulk has his own task...
	task, terr := c.AsyncTask(resp.Payload.Task[0])
	if terr != nil {
		return nil, goof.Newf("->CreateVolumeSnapshot(), %v", terr)
	}

	return c.ShowVolume(task.Resource.ID)
}

// CreateSnapshotFullCopy Use gocoprhd to get the list of volumes
// <- VolumeCreateFromSnapshot
func (c *CoprHDClient) CreateSnapshotFullCopy(snapid string, volname string) (*apimodels.Volume, error) {

	// Construct Request Parameters
	body := &apimodels.CreateSnapshotFullCopy{
		Count:          1,
		Name:           volname,
		CreateInactive: false,
		Type:           "rp",
	}

	// Create Object to Request
	createSnapshotFullCopyParams := apiblock.NewCreateSnapshotFullCopyParams().WithID(snapid).WithBody(body)

	//use any function to do REST operations
	resp, err := c.block.CreateSnapshotFullCopy(createSnapshotFullCopyParams, c.authInfo)
	log.Infof("->CreateSnapshotFullCopy(): %#v", resp.Payload)
	if err != nil {
		return nil, goof.Newf("->CreateSnapshotFullCopy(), %v", err)
	}

	// Monitor the progress bar...
	task, terr := c.AsyncTask(resp.Payload)
	if terr != nil {
		return nil, goof.Newf("->CreateSnapshotFullCopy(), %v", err)
	}

	return c.ShowVolume(task.Resource.ID)
}

//CreateVolumeFullCopy <- VolumeCopy
func (c *CoprHDClient) CreateVolumeFullCopy(volid string, volname string) (*apimodels.Volume, error) {

	// Construct Request Parameters
	body := &apimodels.CreateVolumeFullCopy{
		Count:          1,
		Name:           volname,
		CreateInactive: false,
		Type:           "",
	}

	// Create Object to Request
	createVolumeFullCopyParams := apiblock.NewCreateVolumeFullCopyParams().WithID(volid).WithBody(body)

	//use any function to do REST operations
	resp, err := c.block.CreateVolumeFullCopy(createVolumeFullCopyParams, c.authInfo)
	if err != nil {
		return nil, goof.Newf("->CreateVolumeFullCopy(), %v", err)
	}

	// Monitor the progress bar...
	task, terr := c.AsyncTask(resp.Payload.Task[0])
	if terr != nil {
		return nil, goof.Newf("->CreateVolumeFullCopy(), %v", terr)
	}

	return c.ShowVolume(task.Resource.ID)
}

// CreateVolumeSnapshot <- VolumeSnapshot
func (c *CoprHDClient) CreateVolumeSnapshot(snapname string, volumeid string) (*apimodels.Snapshot, error) {

	// Create the request Parameters
	body := &apimodels.CreateVolumeSnapshot{
		Name:           snapname,
		CreateInactive: true,
		ReadOnly:       false,
	}

	// Create Object to Request
	createVolumeSnapshotParams := apiblock.NewCreateVolumeSnapshotParams().WithID(volumeid).WithBody(body)

	//use any function to do REST operations
	// This method doesn't requires any parameter.
	resp, err := c.block.CreateVolumeSnapshot(createVolumeSnapshotParams, c.authInfo)
	log.Infof("->CreateVolumeSnapshot(): %#v", resp.Payload)
	if err != nil {
		return nil, goof.Newf("->CreateVolumeSnapshot(), %v", err)
	}

	// This task is async so we need to monitor the progress bar...
	task, terr := c.AsyncTask(resp.Payload.Task[0])
	if terr != nil {
		return nil, goof.Newf("->CreateVolumeSnapshot(), %v", terr)
	}

	return c.ShowSnapshot(task.Resource.ID)
}

// DeleteVolume <- VolumeRemove
func (c *CoprHDClient) DeleteVolume(volid string) (*apimodels.Volume, error) {

	// Create Object to Request
	deleteVolumeParams := apiblock.NewDeleteVolumeParams().WithID(volid)

	//use any function to do REST operations
	resp, err := c.block.DeleteVolume(deleteVolumeParams, c.authInfo)
	log.Infof("->DeleteVolume(): %#v", resp.Payload)
	if err != nil {
		return nil, goof.Newf("->DeleteVolume(), %v", err)
	}

	// Monitor the progress bar...
	task, terr := c.AsyncTask(resp.Payload.Task[0])
	if terr != nil {
		return nil, goof.Newf("->DeleteVolume(), %v", terr)
	}

	return c.ShowVolume(task.Resource.ID)
}

// ListVolumeExports ... List Exports of a Volume
func (c *CoprHDClient) ListVolumeExports(volID string) ([]*apimodels.VolumeExports, error) {

	// Create the request
	listVolumeExportsParams := apiblock.NewListVolumeExportsParams().WithID(volID)

	//use any function to do REST operations
	resp, err := c.block.ListVolumeExports(listVolumeExportsParams, authInfo)
	if err != nil {
		return nil, goof.Newf("->ListVolumeExports(), %v", err)
	}

	return resp.Payload.Itl, nil
}

// CreateExport <-> VolumeAttach
func (c *CoprHDClient) CreateExport(name string, hostsID []string, volID string) ([]*apimodels.VolumeExports, error) {

	// Construct Request Parameters
	body := &apimodels.CreateExport{
		Project: c.config.Project(),
		Varray:  c.config.VArray(),
		Name:    name,
		Type:    "Host",
		Hosts:   hostsID,
		Volumes: []*models.CreateExportVolumesItems0{
			&models.CreateExportVolumesItems0{
				ID: volID,
			},
		},
	}

	// Create Object to Request
	createExportParams := apicompute.NewCreateExportParams().WithBody(body)

	//use any function to do REST operations
	resp, err := c.compute.CreateExport(createExportParams, authInfo)
	log.Infof("->CreateExport(): %#v", resp.Payload)
	if err != nil {
		return nil, goof.Newf("->CreateExport(), %v", err)
	}

	// Monitor the progress bar...
	task, terr := c.AsyncTask(resp.Payload.Task[0])
	if terr != nil {
		return nil, goof.Newf("->CreateExport(), %v", terr)
	}

	return c.ListVolumeExports()
}

// DeleteExport <-> VolumeDetach
func (c *CoprHDClient) DeleteExport() (*apimodels.Volume, error) {}

// Snapshots <-> Snapshots
func (c *CoprHDClient) Snapshots() (*apimodels.Volume, error) {

}

// ShowSnapshot <-> SnapshotInspect
func (c *CoprHDClient) ShowSnapshot() (*apimodels.Volume, error) {}

// CreateSnapshotCopy <-> SnapshotCopy
func (c *CoprHDClient) CreateSnapshotCopy() (*apimodels.Snapshot, error) {
	return nil, types.ErrNotImplemented
}

// DeleteSnapshot <-> SnapshotRemove
func (c *CoprHDClient) DeleteSnapshot() (*apimodels.Volume, error) {
}
