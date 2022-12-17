package config

type SloopConfig struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"Kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}
type Global struct {
	Namespace string `yaml:"namespace"`
}
type Set struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
type Components struct {
	Name        string   `yaml:"name"`
	Enabled     bool     `yaml:"enabled"`
	Path        string   `yaml:"path"`
	Namespace   string   `yaml:"namespace"`
	Set         []Set    `yaml:"set"`
	ValuesFiles []string `yaml:"valuesFiles"`
}
type Spec struct {
	Template   string       `yaml:"template"`
	Version    string       `yaml:"version"`
	Global     Global       `yaml:"global"`
	Components []Components `yaml:"components"`
}
