package main

import (
	"net/url"
	"strings"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"fmt"
)

type SiteStats struct {
	Section map[string]SectionStats
}

func NewSiteStats() *SiteStats {
	return &SiteStats{
		make(map[string]SectionStats)}
}

type SectionStats struct {
	Hits       int
	StatusCode map[int]int
}

func (stats *SiteStats) Display(table *tablewriter.Table) {
	// Section stats
	table.SetHeader([]string{"Section", "Hits", "Status Code"})
	for key := range stats.Section {
		value, _ := stats.Section[key]
		statusCodes := ""
		for code := range value.StatusCode {
			statusCodes += fmt.Sprintf("%d: %d, ", code,  value.StatusCode[code])
		}
		table.Append([]string{key, strconv.Itoa(value.Hits), statusCodes[:len(statusCodes)-2]}) // remove trailing ","
	}
	table.Render()
	table.ClearRows()
}

func (stats *SiteStats) Update(logEntries []*CommonLogEntry) {
	for _, entry := range logEntries {
		// Section stats
		stats.updateSectionStats(entry)
	}
}

func (stats *SiteStats) updateSectionStats(entry *CommonLogEntry) {
	section := parseSection(entry.Request.URL)
	statusCode := entry.StatusCode
	if value, ok := stats.Section[section]; ok == true {
		value.Hits++
		value.StatusCode[statusCode]++
	} else {
		stats.Section[section] = SectionStats{Hits: 1, StatusCode: map[int]int{statusCode: 1}}
	}
}

func parseSection(url *url.URL) string {
	if pathArray := strings.Split(url.EscapedPath(), "/"); len(pathArray) >= 3 {
		return strings.TrimSpace(pathArray[1])
	} else {
		return "home"
	}
}
