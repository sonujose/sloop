package console

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func newTableInstance() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(getDefaultTableStyle())
	return t
}

func getDefaultTableStyle() table.Style {

	l1 := table.Style{
		Name:    "StyleColoredDefaultPlaneTable",
		Box:     table.StyleBoxDefault,
		Color:   table.ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		HTML:    table.DefaultHTMLOptions,
		Options: table.OptionsNoBordersAndSeparators,
		Title:   table.TitleOptionsBlackOnBlue,
	}

	return l1
}
