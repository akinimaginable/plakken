package constant

import "time"

const (
	HTTPTimeout          = 3 * time.Second
	ExpirationCurlCreate = 604800 * time.Second // Second in one week
)
