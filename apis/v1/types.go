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
type HelmSources struct {
	Name      string   `yaml:"name"`
	URL       string   `yaml:"url,omitempty"`
	Remote    bool     `yaml:"remote"`
	Namespace string   `yaml:"namespace"`
	Values    []string `yaml:"values"`
	Path      string   `yaml:"path,omitempty"`
}
type Sources struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Path      string   `yaml:"path"`
	Namespace string   `yaml:"namespace"`
	Values    []string `yaml:"values"`
}
type Spec struct {
	Template    string        `yaml:"template"`
	Version     string        `yaml:"version"`
	Global      Global        `yaml:"global"`
	HelmSources []HelmSources `yaml:"helmSources"`
	Sources     []Sources     `yaml:"sources"`
}
