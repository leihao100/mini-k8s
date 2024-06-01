package status

import "encoding/json"

type DNSStatus struct {
	Configured bool `json:"configured"`
}

func (d *DNSStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DNSStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}
