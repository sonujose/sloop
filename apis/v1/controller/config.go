package controller

type SloopControllerConfig struct {
	Name   string `json:"name"`
	Config Config `json:"config"`
	Status Status `json:"status"`
}
type TemplateFile struct {
	FileName     string `json:"fileName"`
	ManifestYaml string `json:"manifest_yaml"`
}
type Component struct {
	Name          string         `json:"name"`
	Namespace     string         `json:"namespace"`
	Path          string         `json:"path"`
	TemplateFiles []TemplateFile `json:"templateFiles"`
}
type Config struct {
	Components           []Component `json:"components"`
	ConsolidatedManifest string      `json:"consolidated_manifest"`
}
type Status struct {
	DeployedOn   string `json:"deployed_on" yaml:"updated"`
	SyncRevision int    `json:"sync_revision" yaml:"revision"`
	Version      string `json:"version" yaml:"version"`
	Components   int    `json:"components" yaml:"components"`
	Name         string `json:"name" yaml:"name"`
}
