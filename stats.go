package main

import ()

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

}


