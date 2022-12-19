package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sonujose/sloop/pkg/core/history"
)

func PrintSyncHistory(hs []history.SyncHistoryObj) {

	t := newTableInstance()
	t.AppendHeader(table.Row{"PACKAGE", "REV", "STATUS", "COMPONENTS", "CREATED", "VERSION"})

	for _, j := range hs {
		updatedTime := j.Updated.Format(time.RFC1123)
		t.AppendRow([]interface{}{j.Package, j.Revision, j.Status, getComponentsStatus(j.Components), updatedTime, j.Version})
	}

	t.Render()
}

func getComponentsStatus(cs string) string {
	ss := strings.Split(cs, "_")
	var components string
	if len(ss) == 2 {
		components = fmt.Sprintf("%s/%s", ss[0], ss[1])
	} else {
		components = cs
	}
	return components
}
