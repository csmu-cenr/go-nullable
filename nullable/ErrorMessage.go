package nullable

import "encoding/json"

type ErrorMessage struct {
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Attempted string      `json:"attemped,c"`
	Code      int         `json:"code,omitempty"`
	Stack     []string    `json:"stack,omitempty"`
}

func (e ErrorMessage) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
