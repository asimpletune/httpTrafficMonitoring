package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type View struct {
	hitsTable    *tablewriter.Table
	sectionTable *tablewriter.Table
}

func NewView() *View {
	return &View{hitsTable: tablewriter.NewWriter(os.Stdout),
		sectionTable: tablewriter.NewWriter(os.Stdout)}
}

func (v *View) Display(stats *SiteStats, warning bool) {
	v.displayHits(stats, warning)
	v.displaySection(stats, warning)
}

func (v *View) displayHits(stats *SiteStats, warning bool) {
	v.hitsTable.SetHeader([]string{"Total Hits"})
	v.hitsTable.Append([]string{strconv.Itoa(stats.Hits)})
	if warning {
		v.hitsTable.SetFooter([]string{"ALERT: total hits exceeds threshold"})
	} else {
		v.hitsTable.ClearFooter()
	}
	v.hitsTable.Render()
	v.hitsTable.ClearRows()
}

func (v *View) displaySection(stats *SiteStats, warning bool) {
	v.sectionTable.SetHeader([]string{"Section", "Hits", "Status Code"})
	for key := range stats.Section {
		value, _ := stats.Section[key]
		statusCodeStatString := ""
		for code := range value.StatusCode {
			statusCodeStatString += fmt.Sprintf("%d: %d, ", code, value.StatusCode[code])
		}
		v.sectionTable.Append([]string{key, strconv.Itoa(value.Hits), statusCodeStatString[:len(statusCodeStatString)-2]})
	}
	if warning {
		v.sectionTable.SetFooter([]string{"ALERT: total hits exceeds threshold"})
	} else {
		v.sectionTable.ClearFooter()
	}
	v.sectionTable.Render()
	v.sectionTable.ClearRows()
}
