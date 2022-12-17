package console

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sonujose/sloop/pkg/core/history"
)

func PrintSyncHistory(hs []history.SyncHistoryObj) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"REVISION", "UPDATED", "SLOOP PACKAGE", "VERSION"})

	for _, j := range hs {
		t.AppendRow([]interface{}{j.Revision, j.Updated, j.Package, j.Version})
	}
	t.SetStyle(table.StyleBold)
	t.Render()
}
