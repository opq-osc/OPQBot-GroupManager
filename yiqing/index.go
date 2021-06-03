package yiqing

type YiqingRes struct {
	Title    string `json:"title"`
	Time     string `json:"time"`
	IncrTime string `json:"incrTime"`
	logcation Logcation `json:"logcation"`
	Colums  []Colums `json:"colums"`
	MainReport struct{
		Id int `json:"id"`
		Area string `json:"area"`
		Report string `json:"report"`
		Dateline string `json:"dateline"`
		Date int64 `json:"date"`
	} `json:"mainReport"`
	ContryData struct{
		SureCnt string `json:"sure_cnt"`
		SureNewCnt string `json:"sure_new_cnt"`
		RestSureCnt string `json:"rest_sure_cnt"`
		RestSureCntIncr string `json:"rest_sure_cnt_incr"`
		InputCnt string `json:"input_cnt"`
		HiddenCnt string `json:"hidden_cnt"`
		HiddenCntIncr string `json:"hidden_cnt_incr"`
		CureCnt string `json:"cure_cnt"`
		YstCureCnt string `json:"yst_cure_cnt"`
		YstDieCnt	string `json:"yst_die_cnt"`
		YstLikeCnt string `json:"yst_like_cnt"`
		YstSureCnt string `json:"yst_sure_cnt"`
		YstSureHid string `json:"yst_sure_hid"`
	}
}

type Colums struct {
	Title    string `json:"title"`
	List  	[]List `json:"list"`
}

type List struct {
	Current int64 `json:"current"`
	Incr  string `json:"incr"`
}

type  Logcation struct{
	Province string `json:"province"`
	City string `json:"city"`
} 