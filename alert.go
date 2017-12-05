package main

import (
	"time"
)

const defaultCoolOffInterval = time.Duration(120) * time.Second

type Alert struct {
	Threshold       float64
	CheckInterval   time.Duration
	CoolOffInterval time.Duration
	Stats           *SiteStats
}

func NewAlert(threshold float64, interval time.Duration, stats *SiteStats) *Alert {
	return &Alert{
		Threshold:       threshold,
		CheckInterval:   interval,
		CoolOffInterval: defaultCoolOffInterval,
		Stats:           stats}
}

func (a *Alert) Monitor(warn chan bool) {
	previousHitCount := a.Stats.Hits
	for {
		select {
		case <-time.After(a.CheckInterval):
			if !a.ok(previousHitCount, a.CheckInterval) {
				warn <- true
				a.coolOff(a.Stats.Hits)
				warn <- false
			}
		}
		previousHitCount = a.Stats.Hits
	}
}

func (a *Alert) coolOff(oldHitCount int) {
	select {
	case <-time.After(a.CoolOffInterval):
		if a.ok(oldHitCount, a.CoolOffInterval) {
			return
		} else {
			a.coolOff(a.Stats.Hits)
		}
	}
}

func (a *Alert) ok(previousHitCount int, duration time.Duration) bool {
	difference := a.Stats.Hits - previousHitCount
	average := float64(difference) / float64(duration.Seconds())
	return average < a.Threshold
}
