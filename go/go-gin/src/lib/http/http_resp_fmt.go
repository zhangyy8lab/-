package http

// ObjListResp 获取 多条数据 响应结构体
type ObjListResp struct {
	Data  []interface{} `json:"data"`
	Count int64         `json:"count"`
	Error string        `json:"error"`
}

// ObjResp 获取 单条数据 响应结构体
type ObjResp struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// CommonResp 创建/更新/删除 都使用这个响应体
type CommonResp struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Create201Resp 新增数据时 需要返回id
type Create201Resp struct {
	Id    uint64 `json:"id"`
	Error string `json:"error"`
}
