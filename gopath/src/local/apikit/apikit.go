package apikit

// StatusCode is easy check status
type StatusCode string

const (
	// OK is successed
	OK StatusCode = "OK"

	// NG is failed
	NG StatusCode = "NG"
)

type Success struct {
	Status StatusCode  `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type Failure struct {
	Status StatusCode  `json:"status"`
	Errors interface{} `json:"errors,omitempty"`
}

// NewSuccess is create Response for success
func NewSuccess() Success {
	return Success{Status: OK}
}

// NewFailure is create Response for failure
func NewFailure(data interface{}) Failure {
	return Failure{Status: NG, Errors: data}
}
