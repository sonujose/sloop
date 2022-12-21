package console

import (
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sonujose/sloop/pkg/core/consts"
	"github.com/sonujose/sloop/pkg/core/history"
)

func PrintSyncHistory(hs []history.SyncHistory) {

	t := newTableInstance()
	t.AppendHeader(table.Row{"PACKAGE", "REV", "STATUS", "COMPONENTS", "UPDATED", "VERSION"})

	for _, j := range hs {
		t.AppendRow(getRowValues(j))
	}

	t.Render()
}

func getRowValues(hs history.SyncHistory) []interface{} {

	var t text.Colors
	var td = text.Colors{text.Color(text.Default)}
	if hs.Status == consts.StatusRegistered {
		t = text.Colors{text.FgHiYellow}
	} else if hs.Status == consts.StatusSynced {
		t = text.Colors{text.FgGreen}
	} else {
		t = text.Colors{text.Color(text.Default)}
	}

	return []interface{}{
		td.Sprint(hs.Package),
		td.Sprint(hs.Revision),
		t.Sprint(hs.Status),
		td.Sprint(getComponentsStatus(hs.Components)),
		td.Sprint(hs.Updated.Format(time.RFC1123)),
		td.Sprint(hs.Version),
	}
}

func getComponentsStatus(cs string) string {
	ss := strings.Split(cs, "_")
	t := text.Colors{text.Color(text.Default)}
	var components string
	if len(ss) == 2 {
		components = t.Sprintf("%s/%s", ss[0], ss[1])
	} else {
		components = cs
	}
	return components
}
