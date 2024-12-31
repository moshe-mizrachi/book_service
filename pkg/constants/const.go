package constants

import "time"

type Function int

const (
	DoCreateIndex Function = 0
	DoUpdateIndex Function = 1
	DoDeleteIndex Function = 2
)

const HighestBookPrice = 9999.0

const (
	MaxIdleConnections        = 50
	MaxIdleConnectionsPerHost = 10
	IdleConnectionTimeout     = 90 * time.Second

	DialTimeout   = 30 * time.Second
	KeepAliveTime = 30 * time.Second

	TLSHandshakeTimeout   = 10 * time.Second
	ExpectContinueTimeout = 1 * time.Second
)

const ActionRoute = "/activity"

const ValidatedAccess = "validated"

const (
	FlushSize         = 100
	FlushInterval     = 5 * time.Second
	ActionsChanelSize = 1000
)
