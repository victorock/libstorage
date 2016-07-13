package context

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"

	"github.com/emccode/libstorage/api/types"
)

type lsc struct {
	context.Context
	key             interface{}
	val             interface{}
	req             *http.Request
	right           context.Context
	logger          *log.Logger
	loggerInherited bool
}

func newContext(
	parent context.Context,
	key interface{},
	val interface{},
	req *http.Request,
	right context.Context) *lsc {

	if parent == nil {
		parent = context.Background()
	}

	// figure out who the parent logger instance is. if there is none,
	// reference the log.StandardLogger as the parent.
	var logger *log.Logger
	if ctx, ok := parent.(*lsc); ok {
		logger = ctx.logger
	}
	if logger == nil {
		logger = log.StandardLogger()
	}

	return &lsc{
		Context:         parent,
		key:             key,
		val:             val,
		req:             req,
		right:           right,
		logger:          logger,
		loggerInherited: true,
	}
}

type hasRightSide interface {
	rightSide() context.Context
}

func (ctx *lsc) rightSide() context.Context {
	return ctx.right
}

// New returns a new context with the provided parent.
func New(parent context.Context) types.Context {
	return newContext(parent, nil, nil, nil, nil)
}

// Background returns a new context with logging capabilities.
func Background() types.Context {
	return New(nil)
}

// WithRequestRoute returns a new context with the injected *http.Request
// and Route.
func WithRequestRoute(
	parent context.Context,
	req *http.Request,
	route types.Route) types.Context {

	return newContext(parent, RouteKey, route, req, nil)
}

// WithStorageService returns a new context with the StorageService as the
// value and attempts to assign the service's associated InstanceID and
// LocalDevices (by way of the service's StorageDriver) to the context as well.
func WithStorageService(
	parent context.Context, service types.StorageService) types.Context {

	driverName := strings.ToLower(service.Driver().Name())

	// set the service's InstanceID if present
	if iidm, ok := parent.Value(AllInstanceIDsKey).(types.InstanceIDMap); ok {
		if iid, ok := iidm[driverName]; ok {
			parent = newContext(parent, InstanceIDKey, iid, nil, nil)
		}
	}

	// set the service's LocalDevices if present
	if ldm, ok := parent.Value(AllLocalDevicesKey).(types.LocalDevicesMap); ok {
		if ld, ok := ldm[driverName]; ok {
			parent = newContext(parent, LocalDevicesKey, ld, nil, nil)
		}
	}

	return newContext(parent, ServiceKey, service, nil, nil)
}

// WithValue returns a copy of parent in which the value associated with
// key is val.
func WithValue(ctx context.Context, key, val interface{}) types.Context {
	if customKeyID, ok := isCustomKey(key); ok {
		key = customKeyID
	}
	return newContext(ctx, key, val, nil, nil)
}

func (ctx *lsc) WithValue(key, val interface{}) types.Context {
	return WithValue(ctx, key, val)
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.  Successive calls to Value with
// the same key returns the same result.
func Value(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

func (ctx *lsc) Value(key interface{}) interface{} {

	if key == LoggerKey {
		return ctx.logger
	}

	if customKeyID, ok := isCustomKey(key); ok {
		key = customKeyID
	}

	if key == HTTPRequestKey {
		return ctx.req
	}

	if ctx.req != nil {
		if req, ok := ctx.Context.Value(HTTPRequestKey).(*http.Request); ok {
			if val, ok := gcontext.GetOk(req, key); ok {
				return val
			}
		}
	}

	if ctx.key == key {
		return ctx.val
	}

	if val := ctx.Context.Value(key); val != nil {
		return val
	}

	if ctx.right != nil {
		return ctx.right.Value(key)
	}

	return nil
}

type hasName interface {
	Name() string
}

type hasID interface {
	ID() string
}

func stringValue(ctx context.Context, key interface{}) (string, bool) {
	switch tv := ctx.Value(key).(type) {
	case string:
		return tv, true
	case *string:
		return *tv, true
	case hasName:
		return tv.Name(), true
	case hasID:
		return tv.ID(), true
	case fmt.Stringer:
		return tv.String(), true
	default:
		return "", false
	}
}

// Join joins this context with another, such that value lookups will first
// first check the current context, and if no such value exist, a lookup
// will be performed against the right side.
func Join(left types.Context, right context.Context) types.Context {

	if left == nil {
		return nil
	}

	if right == nil {
		return left
	}

	if left == right {
		return left
	}

	return newContext(left, nil, nil, nil, right)
}
func (ctx *lsc) Join(right context.Context) types.Context {
	return Join(ctx, right)
}

// SetLogLevel sets the context's log level.
func SetLogLevel(ctx context.Context, lvl log.Level) {
	if logCtx, ok := ctx.(*lsc); ok {
		if lvl == logCtx.logger.Level {
			return
		}
		if logCtx.loggerInherited {
			parentLogger := logCtx.logger
			logCtx.logger = &log.Logger{
				Formatter: parentLogger.Formatter,
				Out:       parentLogger.Out,
				Hooks:     parentLogger.Hooks,
				Level:     lvl,
			}
			logCtx.loggerInherited = false
		}
	}
}

// GetLogLevel gets the context's log level.
func GetLogLevel(ctx context.Context) (log.Level, bool) {
	if logCtx, ok := ctx.(*lsc); ok {
		return logCtx.logger.Level, true
	}
	return 0, false
}

// InstanceID returns the context's InstanceID.  This value is valid on both
// the client and the server.
func InstanceID(ctx context.Context) (*types.InstanceID, bool) {
	v, ok := ctx.Value(InstanceIDKey).(*types.InstanceID)
	return v, ok
}

// MustInstanceID returns the context's InstanceID and panics if it does not
// exist and/or cannot be type cast.
func MustInstanceID(ctx context.Context) *types.InstanceID {
	return ctx.Value(InstanceIDKey).(*types.InstanceID)
}

// LocalDevices returns the context's local devices map.  This value is valid
// on both the client and the server.
func LocalDevices(ctx context.Context) (*types.LocalDevices, bool) {
	v, ok := ctx.Value(LocalDevicesKey).(*types.LocalDevices)
	return v, ok
}

// Transaction returns the context's Transaction. This value is valid on both
// the client and the server.
func Transaction(ctx context.Context) (*types.Transaction, bool) {
	v, ok := ctx.Value(TransactionKey).(*types.Transaction)
	return v, ok
}

// MustTransaction returns the context's Transaction and will panic if the
// value is missing or not castable.
func MustTransaction(ctx context.Context) *types.Transaction {
	return ctx.Value(TransactionKey).(*types.Transaction)
}

// RequireTX ensures a context has a transaction, and if it doesn't creates a
// new one.
func RequireTX(ctx context.Context) types.Context {
	if _, ok := Transaction(ctx); ok {
		return newContext(ctx, nil, nil, nil, nil)
	}

	tx, err := types.NewTransaction()
	if err != nil {
		panic(err)
	}

	return newContext(ctx, TransactionKey, tx, nil, nil)
}

// Client returns the context's Client. This value is only valid for
// contexts created on the client.
func Client(ctx context.Context) (types.Client, bool) {
	v, ok := ctx.Value(ClientKey).(types.Client)
	return v, ok
}

// MustClient returns the context's Client and panics if the client is not
// available and/or not castable..
func MustClient(ctx context.Context) types.Client {
	return ctx.Value(ClientKey).(types.Client)
}

// Profile returns the context's profile. This value is only valid for
// contexts created on the server after the profile handler has been executred.
func Profile(ctx context.Context) (string, bool) {
	return stringValue(ctx, ProfileKey)
}

// Route returns the context's route. This value is only valid for contexts
// created on the server after a mux has received an incoming HTTP request.
// Any part of the libStorage workflow after that, including the handlers,
// routers, and storage drivers, should have access to the Route.
func Route(ctx context.Context) (types.Route, bool) {
	v, ok := ctx.Value(RouteKey).(types.Route)
	return v, ok
}

// Server returns the context's server name. This value is valid on both the
// client and the server.
func Server(ctx context.Context) (string, bool) {
	return stringValue(ctx, ServerKey)
}

// Service returns the context's storage service. This value is valid only for
// contexts created on the server. The value is only available after the
// service has been injected as part of the ServiceValidator handler or by
// one of the routes that injects it manually such as Volumes or Snapshots.
func Service(ctx context.Context) (types.StorageService, bool) {
	v, ok := ctx.Value(ServiceKey).(types.StorageService)
	if !ok {
		v, ok = ctx.Value(StorageServiceKey).(types.StorageService)
	}
	return v, ok
}

// MustService returns the context's StorageService  and panics if it does not
// exist and/or cannot be type cast.
func MustService(ctx context.Context) types.StorageService {
	v, _ := Service(ctx)
	return v.(types.StorageService)
}

// ServiceName returns the context's service name. This value is valid for
// contexts created on both the client and the server. On the server this
// value is subject to the same restrictions as listed in the Service function.
func ServiceName(ctx context.Context) (string, bool) {
	v, ok := stringValue(ctx, ServiceKey)
	if !ok {
		v, ok = stringValue(ctx, StorageServiceKey)
	}
	return v, ok
}

// Driver returns the context's storage driver. This value is valid only
// on the server and subject to the same restrictions as listed in the Service
// function.
func Driver(ctx context.Context) (types.StorageDriver, bool) {
	v, ok := ctx.Value(DriverKey).(types.StorageDriver)
	if ok {
		return v, ok
	}

	w, ok := ctx.Value(StorageServiceKey).(types.StorageService)
	if ok {
		if d := w.Driver(); d != nil {
			return d, true
		}
	}

	return nil, false
}
