package v1

type SloopConfigmap struct {
	Name      string      `json:"name"`
	Config    Config      `json:"config"`
	Manifests []Manifests `json:"manifests"`
	Info      Info        `json:"info"`
}
type ManifestsList struct {
	Name         string `json:"name"`
	ManifestYaml string `json:"manifest_yaml"`
}
type Components struct {
	Name          string          `json:"name"`
	Namespace     string          `json:"namespace"`
	Path          string          `json:"path"`
	ManifestsList []ManifestsList `json:"manifests_list"`
}
type Config struct {
	Components []Components `json:"components"`
}
type Manifests struct {
	Component    string `json:"component"`
	ManifestYaml string `json:"manifest_yaml"`
}
type Info struct {
	DeployedAt string `json:"deployed_at"`
}
