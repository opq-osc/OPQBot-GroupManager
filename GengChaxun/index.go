package GengChaxun

type GenChaxunRes struct {
	Data 	[]DataRes `json:"data"`
}

type DataRes struct {
	Definitions []Items `json:"definitions"`
	Tags []TageItem `json:"tags"`
}

type Items struct {
	Content string `json:"content"`
	Plaintext string `json:"plaintext"`
	Images []ImageItem `json:"images"`
}

type ImageItem struct {
	Full  FullItem `json:"full"`
}

type FullItem struct {
	Path string `json:"path"`
}
type TageItem struct {
	Name string `json:"name"`
}