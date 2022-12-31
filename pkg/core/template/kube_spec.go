package template

import (
	"bytes"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yaml2 "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// addKubeSpecToManifest - Injects Kubernetes spec to the manifests - labels and namespace
func (t *SloopTemplate) injectKubeSpecToManifest(manifestYamlBlob bytes.Buffer, comp client.Component, log *logrus.Logger) ([]byte, error) {

	manifestBlob := manifestYamlBlob.Bytes()
	kresource := &unstructured.Unstructured{}

	var decUnstructured = yaml2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := decUnstructured.Decode(manifestBlob, nil, kresource); err != nil {
		log.Error(err, "Error parsing manifest blob to client.Object")
	}

	existingResourceLabels := kresource.GetLabels()

	updatedLabels, err := getlabelsForSloopManagedResource(t.sloopPkg.Metadata.Name, comp.Name, existingResourceLabels, log)
	if err != nil {
		return nil, err
	}

	kresource.SetNamespace(t.getNamespaceForResource(comp))
	kresource.SetLabels(updatedLabels)

	result, err := yaml.Marshal(kresource.Object)

	if err != nil {
		log.Error("Error converting kubernetes unstruct resource to manifest")
		return nil, err
	}

	return result, nil
}

func getlabelsForSloopManagedResource(sloopPkg string, component string, existinglables map[string]string, log *logrus.Logger) (map[string]string, error) {

	sloopResourceLabels := map[string]string{
		"app.kubernetes.io/owner": "sloop",
		"sloop.io/component":      component,
		"sloop.io/package":        sloopPkg,
	}

	err := mergo.Merge(&existinglables, sloopResourceLabels, mergo.WithOverride)

	if err != nil {
		log.Error("Error merging existing labels withsloop specific labels")
		return existinglables, err
	}

	return existinglables, nil

}

func (t *SloopTemplate) getNamespaceForResource(comp client.Component) string {

	var templateNamespace string

	if comp.Namespace != "" {
		templateNamespace = comp.Namespace
	} else if t.sloopPkg.Spec.Global.Namespace != "" {
		templateNamespace = t.sloopPkg.Spec.Global.Namespace
	} else if t.sloopPkg.Metadata.Namespace != "" {
		templateNamespace = t.sloopPkg.Metadata.Namespace
	} else {
		templateNamespace = "default"
	}

	return templateNamespace
}
