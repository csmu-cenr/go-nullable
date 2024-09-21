package nullable

import "encoding/json"

type ErrorMessage struct {
	Attempted string      `json:"attemped,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	ErrorNo   int         `json:"errorNo,omitempty"`
	Message   string      `json:"message"`
	Stack     []string    `json:"stack,omitempty"`
}

func (e ErrorMessage) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
