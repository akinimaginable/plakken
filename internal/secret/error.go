package secret

type parseError struct {
	message string
}

func (m *parseError) Error() string {
	return "parseHash: " + m.message
}
