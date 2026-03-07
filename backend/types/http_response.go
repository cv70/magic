package types

type HTTPResponse[T any] struct {
	Code int `json:"code,omitempty"`
	Msg string `json:"msg,omitempty"`
	Data T `json:"data,omitempty"`
}
