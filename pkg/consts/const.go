package consts

import (
	"math"
	"time"
)

// Function Index operations
type Function int

const (
	DoCreateIndex Function = 0
	DoUpdateIndex Function = 1
	DoDeleteIndex Function = 2
)

// HighestBookPrice Book max price
var HighestBookPrice = math.Inf(1)

// Elasticsearch config
const (
	MaxIdleConnections        = 50
	MaxIdleConnectionsPerHost = 10
	IdleConnectionTimeout     = 90 * time.Second

	DialTimeout   = 30 * time.Second
	KeepAliveTime = 30 * time.Second

	TLSHandshakeTimeout   = 10 * time.Second
	ExpectContinueTimeout = 1 * time.Second
	WorkersNumber         = 10
)

// ActionRoute routes
const ActionRoute = "activity"

// ValidatedAccess Validations
const ValidatedAccess = "validated"

// Redis config
const (
	FlushSize         = 100
	FlushInterval     = 5 * time.Second
	ActionsChanelSize = 1000
)
