package pkg

import "encoding/json"

type Error struct {
	status int
	error  string
}

func (e Error) StatusCode() int {
	return e.status
}

func (e Error) Error() string {
	return e.error
}

func (e Error) MarshalJSON() ([]byte, error) {
	if e.status == 400 {
		return json.Marshal(map[string]string{"err": e.Error()})
	}
	return nil, nil
}

func makeError(status int, message string) error {
	return Error{status, message}
}
