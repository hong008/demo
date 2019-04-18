package module

import "net/http"

type Data interface {
	Valid() bool //判断数据有效性
}

/*数据请求*/
type Request struct {
	httpReq *http.Request //http请求
	depth   uint32        //请求深度
}

func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{
		httpReq: httpReq,
		depth:   depth,
	}
}

func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() uint32 {
	return req.depth
}

//判断爬虫请求是否有效
func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

/*数据响应*/
type Response struct {
	httpResp *http.Response //http响应
	depth    uint32         //响应的深度，对应请求的深度
}

func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{
		httpResp: httpResp,
		depth:    depth,
	}
}

func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

/*数据条目*/
type Item map[string]interface{}

func (i Item) Valid() bool {
	return i != nil
}
