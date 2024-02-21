package database

type PingError struct{}

func (m *PingError) Error() string {
	return "Connection to redis not work"
}
