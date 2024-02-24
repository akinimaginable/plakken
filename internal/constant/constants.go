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
)
