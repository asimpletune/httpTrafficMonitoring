package main

import (
	"net/url"
	"strings"
)

type SiteStats struct {
	Hits    int
	Section map[string]SectionStats
}

func NewSiteStats() *SiteStats {
	return &SiteStats{
		0,
		make(map[string]SectionStats)}
}

type SectionStats struct {
	Hits       int
	StatusCode map[int]int
}

func (stats *SiteStats) Update(logEntries []*CommonLogEntry) {
	for _, entry := range logEntries {
		// Overall hit stats
		stats.Hits++
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
