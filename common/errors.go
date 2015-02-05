package common

type ServerError struct {
	Err error
}

func (se ServerError) Error() string {
	return se.Err.Error()
}

func (se ServerError) String() string {
	return se.Err.Error()
}
