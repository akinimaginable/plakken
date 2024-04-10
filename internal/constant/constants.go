package constant

import "time"

const (
	HTTPTimeout          = 3 * time.Second
	ExpirationCurlCreate = 604800 * time.Second // Second in one week
	TokenLength          = 32
	ArgonSaltSize        = 16
	ArgonMemory          = 64 * 1024
	ArgonThreads         = 4
	ArgonKeyLength       = 32
	ArgonIterations      = 2
	MaxURLLength         = 255
	SecondsInDay         = 86400
	SecondsInHour        = 3600
	SecondsInMinute      = 60
)
