package err

type Error string
func (e Error) Error() string {
	return string(e)
}

const (
	ErrWsParseFailed = Error("parse websocket message failed")
	ErrWsInvalidType = Error("invalid websocket message type")
)
