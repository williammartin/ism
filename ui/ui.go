package ui

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/lunixbochs/vtclean"
	runewidth "github.com/mattn/go-runewidth"
)

const DefaultTableSpacePadding = 3

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

func (ui *UI) DisplayTable(table [][]string) {
	if len(table) == 0 {
		return
	}

	for i, str := range table[0] {
		style := color.New(color.Bold)
		table[0][i] = style.Sprint(str)
	}

	var columnPadding []int

	rows := len(table)
	columns := len(table[0])
	for col := 0; col < columns; col++ {
		var max int
		for row := 0; row < rows; row++ {
			if strLen := wordSize(table[row][col]); max < strLen {
				max = strLen
			}
		}
		columnPadding = append(columnPadding, max+DefaultTableSpacePadding)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			data := table[row][col]
			var addedPadding int
			if col+1 != columns {
				addedPadding = columnPadding[col] - wordSize(data)
			}
			fmt.Fprintf(ui.Out, "%s%s", data, strings.Repeat(" ", addedPadding))
		}
		fmt.Fprintf(ui.Out, "\n")
	}
}

func wordSize(str string) int {
	cleanStr := vtclean.Clean(str, false)
	return runewidth.StringWidth(cleanStr)
}
