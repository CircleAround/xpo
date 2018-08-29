package apikit

// StatusCode is easy check status
type StatusCode string

const (
	// OK is successed
	OK StatusCode = "OK"

	// NG is failed
	NG StatusCode = "NG"
)

// ResponseWrapper is API Response Structure
type ResponseWrapper struct {
	Status StatusCode  `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// NewSuccess is create Response for success
func NewSuccess() ResponseWrapper {
	return ResponseWrapper{Status: OK}
}

// NewFailure is create Response for failure
func NewFailure(data interface{}) ResponseWrapper {
	return ResponseWrapper{Status: NG, Data: data}
}


