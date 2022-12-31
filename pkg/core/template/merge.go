package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/sonujose/sloop/apis/v1/client"
)

// generateValuesForOverrides - Generate ValuesFiles based on the overrides specified in the sloop config
func generateValuesForOverrides(keys []string, val string) map[string]interface{} {

	finalElement := map[string]interface{}{
		keys[len(keys)-1]: val,
	}

	for i := len(keys) - 2; i > 0; i-- {
		newElement := map[string]interface{}{
			keys[i]: finalElement,
		}
		finalElement = newElement
	}

	return finalElement
}

// templateExecutor - Template Executor for go templating
func templateExecutor(t []byte, vals map[string]interface{}, out io.Writer) error {
	tt, err := template.New("_").Parse(string(t))
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}

// appendTemplateManifestWithRef - Append the Metaline for the template files.
func appendTemplateManifestWithRef(templateManifest []byte, comp client.Component, templateFile string) bytes.Buffer {

	var packageTemplateManifest bytes.Buffer

	// Write the whole template manifest for all components
	packageTemplateManifest.Write([]byte("---\n"))
	packageTemplateManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s, Namespace: %s\n", comp.Name, templateFile, comp.Namespace)))
	packageTemplateManifest.Write(templateManifest)

	return packageTemplateManifest
}
