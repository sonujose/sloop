package template

import (
	"io"
	"text/template"
)

func tpl(t string, vals map[string]interface{}, out io.Writer) error {
	tt, err := template.New("_").Parse(t)
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}
