package services

import (
	"net/http"
	"time"
)

type HeartbeatOperations struct {
	RestOperations
}

func NewHeartbeatOperations() *HeartbeatOperations {
	hbo := HeartbeatOperations{}
	return &hbo
}

func (c *HeartbeatOperations) GetHeartbeatOperation() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		c.Heartbeat(res, req)
	}
}

func (c *HeartbeatOperations) Heartbeat(res http.ResponseWriter, req *http.Request) {
	c.SendResult(res, req, time.Now(), nil)
}
