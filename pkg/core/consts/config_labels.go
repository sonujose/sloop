package consts

import "encoding/json"

type SloopConfigSecretLabels struct {
	Components string `json:"components"`
	Owner      string `json:"owner" default:"sloop"`
	Package    string `json:"package"`
	Revision   string `json:"revision"`
	Status     string `json:"status" default:"registered"`
	Updated    string `json:"updated"`
	Version    string `json:"version"`
}

func (l *SloopConfigSecretLabels) GetConfigLabels() map[string]string {
	d1, _ := json.Marshal(l)
	var labls map[string]string

	json.Unmarshal(d1, &labls)

	return labls
}
