package history

import (
	"github.com/sonujose/sloop/pkg/core/consts"
	klabel "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	// PackageHistoryFilterKey : List history of a specfic package
	PackageHistoryFilterKey string = "package"

	// AllPackagesHistoryFilterKey : List history of all Packages in the sloop namespace
	AllPackagesHistoryFilterKey string = "allPackages"

	//RegisteredPackagesHistoryFilterKey : List history of package revisions with registered status
	RegisteredPackagesHistoryFilterKey string = "packgeRegisteredRevisions"
)

// GetPackageHistoryLabelSelectors - Provides label selectors for getting package history for different purpose
// Package History purpose Filters
// 		- PackageHistoryFilterKey
// 		- AllPackagesHistoryFilterKey
// 		- RegisteredPackagesHistoryFilterKey
func GetPackageHistoryLabelSelectors(purpose string, packageName string) klabel.Selector {

	var labels map[string]string
	switch purpose {

	case PackageHistoryFilterKey:
		labels = map[string]string{"owner": consts.ToolName, "package": packageName}
		break
	case AllPackagesHistoryFilterKey:
		labels = map[string]string{"owner": consts.ToolName}
		break
	case RegisteredPackagesHistoryFilterKey:
		labels = map[string]string{"owner": consts.ToolName, "package": packageName, "status": consts.StatusRegistered}
		break
	default:
		labels = map[string]string{"owner": consts.ToolName}
		break
	}

	return getKubeLabelSelectorFromMaps(labels)
}

func getKubeLabelSelectorFromMaps(labelSelectors map[string]string) klabel.Selector {

	lbsel := klabel.NewSelector()

	for key, val := range labelSelectors {
		req, err := klabel.NewRequirement(key, selection.Operator(selection.Equals), []string{val})
		if err != nil {
			continue
		}
		lbsel = lbsel.Add(*req)
	}

	return lbsel
}
