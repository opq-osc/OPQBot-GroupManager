package GengChaxun

type GenChaxunRes struct {
	Data 	[]DataRes `json:"data"`
}

type DataRes struct {
	Definitions []Items `json:"definitions"`
}

type Items struct {
	Content string `json:"content"`
	Images []ImageItem `json:"images"`
}

type ImageItem struct {
	Full  FullItem `json:"full"`
}

type FullItem struct {
	Path string `json:"path"`
}