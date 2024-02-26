package plak

type deletePlakError struct {
	name string
	err  error
}

func (m *deletePlakError) Error() string {
	return "Cannot delete: " + m.name + " : " + m.err.Error()
}

type createError struct {
	message string
}

func (m *createError) Error() string {
	return "create: cannot create plak: " + m.message
}
