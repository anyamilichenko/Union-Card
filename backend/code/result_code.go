package code

type ResultCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (resultCode *ResultCode) SetMessage(message string) *ResultCode {
	resultCode.Message = message
	return resultCode
}
