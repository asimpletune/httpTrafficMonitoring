package main

import (
	"net/url"
	"strings"
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

func (stats *SiteStats) Update(logEntries []*CommonLogEntry) {
	for _, entry := range logEntries {
		// Section stats
		stats.updateSectionStats(entry)
	}
}

func (stats *SiteStats) updateSectionStats(entry *CommonLogEntry) {
	section := parseSection(entry.Request.URL)
	if value, ok := stats.Section[section]; ok == true {
		value.Hits++
		value.StatusCode[entry.StatusCode]++
	}
}

func parseSection(url *url.URL) string {
	if pathArray := strings.Split(url.EscapedPath(), "/"); len(pathArray) >= 3 {
		return strings.TrimSpace(pathArray[1])
	} else {
		return "home"
	}
}
