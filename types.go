package gsclient

import "encoding/json"

type serverHardwareProfile struct {
	string
}

//MarshalJSON custom marshal for serverHardwareProfile
func (s serverHardwareProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}