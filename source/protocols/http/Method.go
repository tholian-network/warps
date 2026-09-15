package http

type Method string

const (
	MethodConnect Method = "CONNECT"
	MethodDelete  Method = "DELETE"
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodOptions Method = "OPTIONS"
	MethodPatch   Method = "PATCH"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodTrace   Method = "TRACE"
)

func (method Method) String() string {

	switch method {
	case MethodConnect:
		return "CONNECT"
	case MethodDelete:
		return "DELETE"
	case MethodGet:
		return "GET"
	case MethodHead:
		return "HEAD"
	case MethodOptions:
		return "OPTIONS"
	case MethodPatch:
		return "PATCH"
	case MethodPost:
		return "POST"
	case MethodPut:
		return "PUT"
	case MethodTrace:
		return "TRACE"
	default:
		return ""
	}

}
