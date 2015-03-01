package v0

// golang net/http does not support client error codes defined outside
// of RFC 2616 so we define them here.
const (
	httpUnprocessableEntity = 422
)
