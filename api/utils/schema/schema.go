package schema

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cesanta/ucl"
	"github.com/cesanta/validate-json/schema"

	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/api/types/context"
	httptypes "github.com/emccode/libstorage/api/types/http"
)

const (
	jsonSchemaID = "https://github.com/emccode/libstorage"
)

var (
	jsonSchema = []byte(JSONSchema)

	// VolumeSchema is the JSON schema for the Volume resource.
	VolumeSchema = buildSchemaVar("volume")

	// VolumeAttachmentSchema is the JSON schema for the VolumeAttachment
	// resource.
	VolumeAttachmentSchema = buildSchemaVar("volumeAttachment")

	// ServiceVolumeMapSchema is the JSON schema for the ServiceVolumeMap
	// resource.
	ServiceVolumeMapSchema = buildSchemaVar("serviceVolumeMap")

	// ServiceSnapshotMapSchema is the JSON schema for the ServiceSnapshotMap
	// resource.
	ServiceSnapshotMapSchema = buildSchemaVar("serviceSnapshotMap")

	// VolumeMapSchema is the JSON schema for the VolumeMap resource.
	VolumeMapSchema = buildSchemaVar("volumeMap")

	// SnapshotMapSchema is the JSON schema for the SnapshotMap resource.
	SnapshotMapSchema = buildSchemaVar("snapshotMap")

	// SnapshotSchema is the JSON schema for the Snapshot resource.
	SnapshotSchema = buildSchemaVar("snapshot")

	// ServiceInfoSchema is the JSON schema for the ServiceInfo resource.
	ServiceInfoSchema = buildSchemaVar("serviceInfo")

	// ServiceInfoMapSchema is the JSON schemea for a map[string]*ServiceInfo.
	ServiceInfoMapSchema = buildSchemaVar("serviceInfoMap")

	// DriverInfoSchema is the JSON schema for the DriverInfo resource.
	DriverInfoSchema = buildSchemaVar("driverInfo")

	// ExecutorInfoSchema is the JSON schema for the ExecutorInfo resource.
	ExecutorInfoSchema = buildSchemaVar("executorInfo")

	// VolumeCreateRequestSchema is the JSON schema for a Volume creation
	// request.
	VolumeCreateRequestSchema = buildSchemaVar("volumeCreateRequest")

	// VolumeCopyRequestSchema is the JSON schema for a Volume copy
	// request.
	VolumeCopyRequestSchema = buildSchemaVar("volumeCopyRequest")

	// VolumeSnapshotRequestSchema is the JSON schema for a Volume snapshot
	// request.
	VolumeSnapshotRequestSchema = buildSchemaVar("volumeSnapshotRequest")

	// VolumeAttachRequestSchema is the JSON schema for a Volume attach
	// request.
	VolumeAttachRequestSchema = buildSchemaVar("volumeAttachRequest")

	// VolumeDetachRequestSchema is the JSON schema for a Volume detach
	// request.
	VolumeDetachRequestSchema = buildSchemaVar("volumeDetachRequest")

	// SnapshotCopyRequestSchema is the JSON schema for a Snapshot copy
	// request.
	SnapshotCopyRequestSchema = buildSchemaVar("snapshotCopyRequest")

	// VolumeCreateFromSnapshotRequestSchema is the JSON schema for a
	// Volume create from Snapshot request.
	VolumeCreateFromSnapshotRequestSchema = buildSchemaVar(
		"volumeCreateFromSnapshotRequest")
)

func buildSchemaVar(name string) []byte {
	return []byte(fmt.Sprintf(`{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "$ref": "%s#/definitions/%s"
}`, jsonSchemaID, name))
}

// ValidateVolume validates a Volume object using the JSON schema. If the
// object is valid no error is returned. The first return value, the object
// marshaled to JSON, is returned whether or not the validation is successful.
func ValidateVolume(v *types.Volume) ([]byte, error) {
	return validateObject(VolumeSchema, v)
}

// ValidateSnapshot validates a Snapshot object using the JSON schema. If the
// object is valid no error is returned. The first return value, the object
// marshaled to JSON, is returned whether or not the validation is successful.
func ValidateSnapshot(v *types.Snapshot) ([]byte, error) {
	return validateObject(SnapshotSchema, v)
}

// ValidateVolumeCreateRequest validates a VolumeCreateRequest object using the
// JSON schema. If the object is valid no error is returned. The first return
// value, the object marshaled to JSON, is returned whether or not the
// validation is successful.
func ValidateVolumeCreateRequest(
	v *httptypes.VolumeCreateRequest) ([]byte, error) {
	return validateObject(VolumeCreateRequestSchema, v)
}

func validateObject(s []byte, o interface{}) (d []byte, e error) {
	if d, e = json.Marshal(o); e != nil {
		return
	}
	if e = Validate(nil, s, d); e != nil {
		return
	}
	return
}

func getSchemaValidator(s []byte) (*schema.Validator, error) {
	volSchema, err := ucl.Parse(bytes.NewReader(s))
	if err != nil {
		return nil, err
	}

	rootSchema, err := ucl.Parse(bytes.NewReader(jsonSchema))
	if err != nil {
		return nil, err
	}

	loader := schema.NewLoader()

	if err := loader.Add(rootSchema); err != nil {
		return nil, err
	}

	validator, err := schema.NewValidator(volSchema, loader)
	if err != nil {
		return nil, err
	}

	return validator, nil
}

// Validate validates the provided data (d) against the provided schema (s).
func Validate(ctx context.Context, s, d []byte) error {

	var l *log.Logger
	if ctx == nil {
		l = log.StandardLogger()
	} else {
		l = ctx.Log()
	}

	l.WithFields(log.Fields{
		"schema": string(s),
		"body":   string(d),
	}).Debug("validating schema")

	validator, err := getSchemaValidator(s)
	if err != nil {
		return err
	}

	if len(d) == 0 {
		d = []byte("{}")
	}

	data, err := ucl.Parse(bytes.NewReader(d))
	if err != nil {
		return err
	}
	return validator.Validate(data)
}