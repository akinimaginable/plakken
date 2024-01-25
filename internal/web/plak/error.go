package plak

type DeletePlakError struct {
	Name string
	Err  error
}

func (m *DeletePlakError) Error() string {
	return "Cannot delete: " + m.Name + " : " + m.Err.Error()
}
