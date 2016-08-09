package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/akutz/goof"
)

// ServerID get something to use as serverid
func ServerID() (string, error) {
	serverid, err := MachineID()
	if err != nil {
		goof.Newf("%v, falling back to HostID()", err)
		return HostID()
	}
	return serverid, nil
}

// MachineID get /etc/machine-id
func MachineID() (string, error) {
	output, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		return "", goof.New("MachineID: Not Supported")
	}
	return string(output), nil
}

// HostID get /etc/machine-id
func HostID() (string, error) {
	output, err := exec.Command("/usr/bin/hostid", "").Output()
	if err != nil {
		return "", goof.Newf("HostID: %v", err)
	}
	return string(output), nil
}

//NormalizeWWN ...
func NormalizeWWN(wwn string) string {
	//Remove prepend 0x if it exists
	wwn = strings.Replace(wwn, "0x", "", 1)
	//Remove whitespaces
	wwn = strings.TrimSpace(wwn)
	//LeftPad...
	wwn = fmt.Sprintf("%016s", wwn)
	// LowerCase...
	wwn = strings.ToLower(wwn)
	return wwn
}

//TwoDotWWN Transform Wwxn in format xx:xx:xx:xx...
func TwoDotWWN(wwn string) string {

	// Split Normalized WWNN
	awwn := strings.Split(NormalizeWWN(wwn), "")

	// Group in pairs
	// < and not <= to not get in off-by-one situation
	var pairs []string
	for i := 0; i < len(awwn)-1; i += 2 {
		pairs = append(pairs, awwn[i]+awwn[i+1])
	}

	// Get pairs and concatenate then with ':' in between
	wwn = strings.Join(pairs, ":")
	return wwn
}

//AddElement ...
// TODO: Create Function to AddElement if non-existent
/*
func AddElement(elements interface{}, element interface{}) interface{} {

	if elements.Type() != element.Type() {
		panic("Elements and Element must be of same Type")
	}

	if HasElement(elements, element) {
		return elements
	}

	elements = reflect.Append(elements, element)
	return elements
}
*/

//HasElement and return bool
/*
func HasElement(element interface{}, elements []interface{}) bool {

	for _, e := range elements {
		if reflect.DeepEqual(e, element) {
			return true
		}
	}
	return false
}
*/
