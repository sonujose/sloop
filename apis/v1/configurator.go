package v1

type SloopConfigurator struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"Kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Global struct {
	Namespace string   `yaml:"namespace"`
	Values    []string `yaml:"values"`
}

type Repositories struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Set struct {
	Name  string `yaml:"name"`
	Value bool   `yaml:"value"`
}

type Releases struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Chart     string `yaml:"chart"`
	Set       []Set  `yaml:"set"`
}

type HelmFile struct {
	Repositories []Repositories `yaml:"repositories"`
	Releases     []Releases     `yaml:"releases"`
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
	HelmFile   HelmFile     `yaml:"helmFile"`
	Components []Components `yaml:"components"`
}
