package gsclient

import "encoding/json"

type serverHardwareProfile struct {
	string
}

//MarshalJSON custom marshal for serverHardwareProfile
func (s serverHardwareProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

type storageType struct {
	string
}

//MarshalJSON custom marshal for storageType
func (s storageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

type ipAddressType struct {
	int
}

//MarshalJSON custom marshal for ipAddressType
func (i ipAddressType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.int)
}
