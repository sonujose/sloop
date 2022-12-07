package v1

type SlooperConfig struct {
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
type Sources struct {
	Name      string   `yaml:"name"`
	Path      string   `yaml:"path"`
	Namespace string   `yaml:"namespace"`
	Values    []string `yaml:"values"`
}
type Spec struct {
	Template string    `yaml:"template"`
	Global   Global    `yaml:"global"`
	Sources  []Sources `yaml:"sources"`
}
