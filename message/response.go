package message

type Response struct {
	ErrMsg *string
}

func NewResponse(msg string) *Response {
	return &Response{ErrMsg: &msg}
}
