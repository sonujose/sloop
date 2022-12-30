package client

type SloopPackage struct {
	APIVersion string   `yaml:"apiVersion" default:"app.sloop.io/v1"`
	Kind       string   `yaml:"Kind" default:"SloopPackage"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace" default:"sloop"`
}
type Global struct {
	Namespace string `yaml:"namespace"`
}
type Set struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
type Component struct {
	Name        string   `yaml:"name"`
	Enabled     bool     `yaml:"enabled"`
	Path        string   `yaml:"path"`
	Namespace   string   `yaml:"namespace"`
	Set         []Set    `yaml:"set"`
	ValuesFiles []string `yaml:"valuesFiles"`
}
type Spec struct {
	Template   string      `yaml:"template"`
	Version    string      `yaml:"version"`
	Global     Global      `yaml:"global"`
	Components []Component `yaml:"components"`
}
