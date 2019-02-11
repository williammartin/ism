package ui

import (
	"io"
	"text/template"
)

type UI struct {
	Out io.Writer
	Err io.Writer
}

func (ui *UI) DisplayText(text string, data ...map[string]interface{}) {
	var keys interface{}
	if len(data) > 0 {
		keys = data[0]
	}

	formattedTemplate := template.Must(template.New("Display Text").Parse(text + "\n"))
	formattedTemplate.Execute(ui.Out, keys)
}
