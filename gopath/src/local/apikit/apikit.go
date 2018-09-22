package apikit

// StatusCode is easy check status
type StatusCode string

const (
	// OK is successed
	OK StatusCode = "OK"

	// NG is failed
	NG StatusCode = "NG"
)

// Success is response OK
type Success struct {
	Status StatusCode  `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// Failure is response NG
type Failure struct {
	Status StatusCode  `json:"status"`
	Error  interface{} `json:"error,omitempty"`
}

// NewSuccess is create Response for success
func NewSuccess() Success {
	return Success{Status: OK}
}

// NewFailure is create Response for failure
func NewFailure(data interface{}) Failure {
	return Failure{Status: NG, Error: data}
}

type Params interface {
	Get(key string) string
	AsInt64(key string) (int64, error)
}
