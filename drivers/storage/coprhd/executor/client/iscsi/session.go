package iscsi

import (
	"io/ioutil"
	"path/filepath"

	"github.com/akutz/goof"
)

const (
	//ISCSISESSIONDIR ...
	ISCSISESSIONDIR = "/sys/class/iscsi_session"
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

//Sessions list sessions from all hosts
// ls -1d /sys/class/iscsi_session/session*
func Sessions() []string {
	sessions, _ := filepath.Glob(ISCSISESSIONDIR + "/session*")
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
	return ISCSISESSIONDIR + "/" + c.SessionID()
}

//Erl ...
func (c *Session) Erl() (string, error) {

	file, err := ioutil.ReadFile(c.BasePath() + "/erl")
	if err != nil {
		return "", goof.Newf("->Session->Erl()%v", err)
	}

	return string(file), nil
}

//RecoveryTmo ...
func (c *Session) RecoveryTmo() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/recovery_tmo")
	if err != nil {
		return "", goof.Newf("->Session->RecoveryTmo()%v", err)
	}

	return string(file), nil
}

//Tpgt ...
func (c *Session) Tpgt() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/tpgt")
	if err != nil {
		return "", goof.Newf("->Session->Tpgt()%v", err)
	}

	return string(file), nil
}

//FirstBurstLen ...
func (c *Session) FirstBurstLen() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/first_burst_len")
	if err != nil {
		return "", goof.Newf("->Session->FirstBurstLen()%v", err)
	}

	return string(file), nil
}

//Creator ...
func (c *Session) Creator() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/creator")
	if err != nil {
		return "", goof.Newf("->Session->Creator()%v", err)
	}

	return string(file), nil
}

//State ...
func (c *Session) State() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/state")
	if err != nil {
		return "", goof.Newf("->Session->State()%v", err)
	}

	return string(file), nil
}

//FastAbort ...
func (c *Session) FastAbort() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/fast_abort")
	if err != nil {
		return "", goof.Newf("->Session->FastAbort()%v", err)
	}

	return string(file), nil
}

//TargetID ...
func (c *Session) TargetID() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/target_id")
	if err != nil {
		return "", goof.Newf("->Session->TargetID()%v", err)
	}

	return string(file), nil
}

//IfaceName ...
func (c *Session) IfaceName() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/ifacename")
	if err != nil {
		return "", goof.Newf("->Session->IfaceName()%v", err)
	}

	return string(file), nil
}

//ImmediateData ...
func (c *Session) ImmediateData() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/immediate_data")
	if err != nil {
		return "", goof.Newf("->Session->ImmediateData()%v", err)
	}

	return string(file), nil
}

//InitialR2T ...
func (c *Session) InitialR2T() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/initial_r2t")
	if err != nil {
		return "", goof.Newf("->Session->InitialR2T()%v", err)
	}

	return string(file), nil
}

//TargetName ...
func (c *Session) TargetName() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/targetname")
	if err != nil {
		return "", goof.Newf("->Session->TargetName()%v", err)
	}

	return string(file), nil
}

//Password ...
func (c *Session) Password() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/password")
	if err != nil {
		return "", goof.Newf("->Session->Password()%v", err)
	}

	return string(file), nil
}

//DataPduInOrder ...
func (c *Session) DataPduInOrder() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/data_pdu_in_order")
	if err != nil {
		return "", goof.Newf("->Session->DataPduInOrder()%v", err)
	}

	return string(file), nil
}

//DataSeqInOrder ...
func (c *Session) DataSeqInOrder() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/data_seq_in_order")
	if err != nil {
		return "", goof.Newf("->Session->DataSeqInOrder()%v", err)
	}

	return string(file), nil
}

//PasswordIn ...
func (c *Session) PasswordIn() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/password_in")
	if err != nil {
		return "", goof.Newf("->Session->PasswordIn()%v", err)
	}

	return string(file), nil
}

//AbortTmo ...
func (c *Session) AbortTmo() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/abort_tmo")
	if err != nil {
		return "", goof.Newf("->Session->AbortTmo()%v", err)
	}

	return string(file), nil
}

//UsernameIn ...
func (c *Session) UsernameIn() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/username_in")
	if err != nil {
		return "", goof.Newf("->Session->UsernameIn()%v", err)
	}

	return string(file), nil
}

//Uevent ...
func (c *Session) Uevent() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/uevent")
	if err != nil {
		return "", goof.Newf("->Session->Uevent()%v", err)
	}

	return string(file), nil
}

//TgtResetTmo ...
func (c *Session) TgtResetTmo() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/tgt_reset_tmo")
	if err != nil {
		return "", goof.Newf("->Session->TgtResetTmo()%v", err)
	}

	return string(file), nil
}

//Username ...
func (c *Session) Username() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/username")
	if err != nil {
		return "", goof.Newf("->Session->Username()%v", err)
	}

	return string(file), nil
}

//MaxBurstLen ...
func (c *Session) MaxBurstLen() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/max_burst_len")
	if err != nil {
		return "", goof.Newf("->Session->MaxBurstLen()%v", err)
	}

	return string(file), nil
}

//LuReseTTmo ...
func (c *Session) LuReseTTmo() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/lu_reset_tmo")
	if err != nil {
		return "", goof.Newf("->Session->LuReseTTmo()%v", err)
	}

	return string(file), nil
}

//MaxOutstandingR2T ...
func (c *Session) MaxOutstandingR2T() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/max_outstanding_r2t")
	if err != nil {
		return "", goof.Newf("->Session->MaxOutstandingR2T()%v", err)
	}

	return string(file), nil
}

//InitiatorName ...
func (c *Session) InitiatorName() (string, error) {
	file, err := ioutil.ReadFile(c.BasePath() + "/initiatorname")
	if err != nil {
		return "", goof.Newf("->Session->InitiatorName()%v", err)
	}

	return string(file), nil
}

//Targets contextual to this session
func (c *Session) Targets() []*Target {
	c.targets = nil
	//target[Host]:[Channel]:[Id]/[Host]:[Channel]:[Id]:[Lun]
	IDs, _ := filepath.Glob(c.BasePath() + "/device/target[0-9]:[0-9]:[0-9]")
	for _, ID := range IDs {
		c.targets = append(c.targets, NewTarget(c, ID))
	}
	return c.targets
}
