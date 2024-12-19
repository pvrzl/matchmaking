package http

type ResponseTitle = string

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type APIResponseArgs struct {
	Status  string
	Message string
	Data    interface{}
	Error   interface{}
	Meta    interface{}
}

func Send(in APIResponseArgs) APIResponse {
	return APIResponse(in)
}
