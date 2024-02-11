package utils

type ParseIntBeforeSeparatorError struct {
	Message string
}

func (m *ParseIntBeforeSeparatorError) Error() string {
	return "parseIntBeforeSeparator: " + m.Message
}

type ParseExpirationError struct {
	Message string
}

func (m *ParseExpirationError) Error() string {
	return "parseIntBeforeSeparator: " + m.Message
}
