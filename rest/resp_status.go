package rest

type HttpStatus struct {
	Value        int
	ReasonPhrase string
}

var OK = HttpStatus{200, "请求成功"}

var Redirect = HttpStatus{302, "请求重定向"}

var BadRequest = HttpStatus{400, "无效的请求"}

var Unauthorized = HttpStatus{401, "token过期"}

var Forbidden = HttpStatus{403, "已被禁止的"}

var NotFound = HttpStatus{404, "没有找到资源"}

var ServerError = HttpStatus{500, "服务器错误"}

var InterfaceError = HttpStatus{501, "接口逻辑错误"}

var BadGateway = HttpStatus{502, "网关错误"}

func GetHttpStatus(value int) (status *HttpStatus) {
	switch value {
	case 200:
		status = &OK
	case 302:
		status = &Redirect
	case 400:
		status = &BadRequest
	case 401:
		status = &Unauthorized
	case 403:
		status = &Forbidden
	case 404:
		status = &NotFound
	case 500:
		status = &ServerError
	case 501:
		status = &InterfaceError
	case 502:
		status = &BadGateway
	default:
		return
	}

	return
}

func GetReason(value int) string {
	switch value {
	case 200:
		return OK.ReasonPhrase
	case 302:
		return Redirect.ReasonPhrase
	case 400:
		return BadRequest.ReasonPhrase
	case 401:
		return Unauthorized.ReasonPhrase
	case 403:
		return Forbidden.ReasonPhrase
	case 404:
		return NotFound.ReasonPhrase
	case 500:
		return ServerError.ReasonPhrase
	case 501:
		return InterfaceError.ReasonPhrase
	case 502:
		return BadGateway.ReasonPhrase
	default:
		return ""
	}
}
