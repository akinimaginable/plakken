package database

type pingError struct{}

func (m *pingError) Error() string {
	return "Connection to redis not work"
}
