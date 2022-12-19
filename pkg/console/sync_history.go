package console

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sonujose/sloop/pkg/core/history"
)

func PrintSyncHistory(hs []history.SyncHistoryObj) {

	t := newTableInstance()
	t.AppendHeader(table.Row{"REVISION", "UPDATED", "SLOOP PACKAGE", "VERSION", "STATUS", "COMPONENTS"})

	for _, j := range hs {
		updatedTime := j.Updated.Format("2006-01-02 15:04:05")
		t.AppendRow([]interface{}{j.Revision, updatedTime, j.Package, j.Version, j.Staus, j.Components})
	}

	t.Render()
}
