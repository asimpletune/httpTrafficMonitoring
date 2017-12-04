package main

import (
	"log"
	"net/http"
	"strconv"
)

type CommonLogEntry struct {
	ClientIp         string
	UserIdentifier   string
	UserId           string
	Date             string
	Request          *http.Request
	StatusCode       int
	ReturnObjectSize string
}

func NewLogEntryFrom(byteArray [][]byte) *CommonLogEntry {
	method := string(byteArray[5][1:])
	url := string(byteArray[6])
	protocol := string(byteArray[7][:len(byteArray[7])])
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Proto = protocol
	statusCode, err := strconv.Atoi(string(byteArray[8]))
	if err != nil {
		log.Fatal(err)
	}
	return &CommonLogEntry{
		string(byteArray[0]),
		string(byteArray[1]),
		string(byteArray[2]),
		string(byteArray[3]) + string(byteArray[4]),
		request,
		statusCode,
		string(byteArray[9]),
	}
}
