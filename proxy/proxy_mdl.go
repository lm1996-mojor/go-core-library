package proxy

type RemoteReqMdl struct {
	Url        string            // 请求地址
	ReqParam   []byte            // 请求参数
	HeadParams map[string]string // 请求头信息
	Method     string            // 请求方式（HTTP 8大请求方式）
	CarryToken bool              // 请求是否需要携带token(默认为false)
	Token      string            // token原文
}

// ParameterlessConstructionOfRemoteReqMdl 无参构造函数
func ParameterlessConstructionOfRemoteReqMdl() RemoteReqMdl {
	return RemoteReqMdl{}
}

// ParametricConstructionOfRemoteReqMdl 带参构造函数
//
// reqParam 请求参数
//
// headParams 请求头信息
//
// url 请求地址
//
// method 请求方式（HTTP 8大请求方式）
//
// carryToken 请求是否需要携带token(默认为false)
//
// token token原文
func ParametricConstructionOfRemoteReqMdl(reqParam []byte, headParams map[string]string, url string, method string, carryToken bool, token string) RemoteReqMdl {
	return RemoteReqMdl{
		url, reqParam, headParams, method, carryToken, token,
	}
}
