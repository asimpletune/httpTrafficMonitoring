package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type Monitor struct{}

func (m *Monitor) GetUpdates(file *os.File) []*CommonLogEntry {
	if recentContent, err := ioutil.ReadAll(file); err == nil {
		logLines := bytes.Split(recentContent, []byte{'\n'})
		return parseLogUpdates(logLines)
	} else {
		log.Fatal(err)
		return nil
	}
}

func parseLogUpdates(byteLines [][]byte) []*CommonLogEntry {
	results := make([]*CommonLogEntry, 0)
	for i := 0; i < len(byteLines); i++ {
		if fields := bytes.Fields(byteLines[i]); len(fields) == 10 {
			logEntry := NewLogEntryFrom(fields)
			results = append(results, logEntry)
		} else {
			break // last line is blank, oftentimes
		}
	}
	return results
}
