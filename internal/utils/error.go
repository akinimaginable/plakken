package utils

type parseIntBeforeSeparatorError struct {
	message string
}

func (m *parseIntBeforeSeparatorError) Error() string {
	return "parseIntBeforeSeparator: " + m.message
}

type ParseExpirationError struct {
	message string
}

func (m *ParseExpirationError) Error() string {
	return "parseIntBeforeSeparator: " + m.message
}
